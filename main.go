package main

import (
    "bufio"
    "bytes"
    "fmt"
    "go/ast"
    "go/parser"
    "go/scanner"
    "go/token"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
    "reflect"
    "strconv"
    "unicode/utf8"
)

type FuncVisitor struct {
    stack *Stack
    buffer []byte
    file *token.File
    root *File
    current token.Position
    lastNode interface{}
    lastPosition []int // shortcut to access last Span/FooterSpan
}

func identToString(ident *ast.Ident) (string) {
  if ident == nil {
    return "<nil>"
  } else {
    return ident.Name
  }
}

func litToString(lit *ast.BasicLit) (string) {
  if lit == nil {
    return "<nil>"
  } else {
    s, err := strconv.Unquote(lit.Value)
    if err != nil {
      return lit.Value
    }
    return s
  }
}

func joinIdentifiers(identifiers []*ast.Ident, c rune) string {
    var buffer bytes.Buffer
    for i := 0; i < len(identifiers); i++ {
      if i > 0 { buffer.WriteRune(c)  }
      buffer.WriteString(identifiers[i].Name)
    }
    return buffer.String()
}

var counter uint = 1
func newid() uint { counter = counter + 1; return counter }

func asNode(node interface{}) *YamlNode {
  switch n := node.(type) {
    case *Terminal:
      y := YamlNode(n)
      return &y
    case *Container:
      y := YamlNode(n)
      return &y
    case *File:
      y := YamlNode(n)
      return &y
  }
  return nil
}

func (v *FuncVisitor) adjustToEOL(endOffset int) int {
  last := len(v.buffer) - 1
  if endOffset >= last { // EOF without EOL
    return endOffset
  }
  c := v.buffer[endOffset]
  for c != '\n' && endOffset < last {
    endOffset++
    c = v.buffer[endOffset]
  }
  return endOffset
}

func (v *FuncVisitor) getName(decl ast.Spec) string {
  switch n := decl.(type) {
    case *ast.ImportSpec:
      return litToString(n.Path)
    case *ast.ValueSpec:
      return joinIdentifiers(n.Names, ',')
    case *ast.TypeSpec:
      return identToString(n.Name)
  }
  return "never happens"
}

func (v *FuncVisitor) getStartOffset(astnode ast.Node) int {
  start := v.file.Position(astnode.Pos());
  previousline := v.current.Line
  newline := start.Line
  if previousline == newline && v.current.Offset < start.Offset {
    // If next token is on the same line, the next starting offset is simply offset + 1
    return v.current.Offset + 1
  } else {
    // Otherwise, move to next EOL lastPosition (already in runes)
    v.lastPosition[1] = v.adjustToEOL(v.lastPosition[1])
    // Check if we reached (or passed) the footer of the parent
    last := asNode(v.lastNode)
    switch p := ((*last).GetParent()).(type) {
      case *Container:
        pFooter := p.GetFooter(0)
        if v.lastPosition[1] >= pFooter { v.lastPosition[1] = pFooter - 1 }
      case *File:
        pFooter := p.GetFooter(0)
        if v.lastPosition[1] >= pFooter { v.lastPosition[1] = pFooter - 1 }
    }
    // Advance current offset
    if v.current.Offset < v.lastPosition[1] {
      v.gotoOffset(v.lastPosition[1])
    }
    return v.current.Offset + 1
  }
}

func (v *FuncVisitor) gotoOffset(offset int) {
  v.current = v.file.Position(v.file.Pos(offset))
}

func (v *FuncVisitor) createTerminal(typeName string, name string, astnode ast.Node) *Terminal {
  end   := v.file.Position(astnode.End());
  endOffset := end.Offset - 1
  startOffset := v.getStartOffset(astnode)
  v.lastPosition = []int{ startOffset, endOffset }
  node := &Terminal {
    id: newid(),
    Type: typeName, Name: name,
    Span: v.lastPosition,
  }
  v.gotoOffset(endOffset)
  return node
}

func (v *FuncVisitor) createSpecialNode(typeName string, name string, startOffset int, endOffset int) {
  v.lastPosition = []int{ startOffset, endOffset }
  node := &Terminal {
    id: newid(),
    Type: typeName, Name: name,
    Span: v.lastPosition,
  }
  v.gotoOffset(endOffset)
  v.root.AddChild(node)
}

func (v *FuncVisitor) createNode(astnode ast.Node, nodename string) interface{} {
    start := v.file.Position(astnode.Pos());
    end   := v.file.Position(astnode.End());
    var node interface{}
    switch n := astnode.(type) {
      case *ast.File:
        packagePos := v.file.Position(n.Name.End())
        v.lastPosition = []int{0, 0}
        v.root = &File {
          id: newid(),
          Type: "file", Name: start.Filename,
          LocationSpan: Location{ Start: []int{start.Line, start.Column}, End: []int{end.Line, end.Column} },
          FooterSpan: []int{ end.Offset, packagePos.Offset - 1 },
        }
        node = v.root
        // We want a package declaration child. Let's treat this as a special case because
        // the original ast.File contains reference to the package offset, but not as a child node
        v.createSpecialNode("PackageDecl", identToString(n.Name), 0, v.root.FooterSpan[1])
      
      case *ast.GenDecl:
        // parenthesized declaration: create container
        if n.Lparen.IsValid() {
          p := v.file.Position(n.Lparen)
          if p.Line > 0 {
            parenOffset := p.Offset
            endOffset := end.Offset - 1
            v.lastPosition = []int{ v.getStartOffset(astnode), parenOffset }
            node = &Container {
              id: newid(),
              Type: n.Tok.String(), Name: n.Tok.String(),
              HeaderSpan: v.lastPosition,
              FooterSpan: []int{ endOffset, v.adjustToEOL(endOffset) },
            }
            v.gotoOffset(parenOffset)
          } else {
            fmt.Println("this should not happen!")
          }
        // If is not parenthesized, create as a terminal
        } else {
          node = v.createTerminal(n.Tok.String(), v.getName(n.Specs[0]), astnode)
        }
      case *ast.ImportSpec:
        var importname string
        if n.Name != nil {
          importname = identToString(n.Name) + ":"
        }
        importname = importname + litToString(n.Path)
        node = v.createTerminal(nodename, importname, astnode)

      case *ast.ValueSpec:
        node = v.createTerminal(nodename, joinIdentifiers(n.Names, ','), astnode)

      case *ast.TypeSpec:
        node = v.createTerminal(nodename, n.Name.Name, astnode)

      case *ast.FuncDecl:   
        node = v.createTerminal(nodename, n.Name.Name, astnode)

      case *ast.Comment:
        node = v.createTerminal(nodename, n.Text, astnode)

      default:
        fmt.Println("Warning! add", nodename, "to isWanted()")
        return nil
    }
    if v.stack.Len() > 0 {
      var p = asNode(v.stack.top.value)
      (*p).AddChild(node)
    }
    
    v.lastNode = node
    return node
}

// We can skip node that we do not want
func (v *FuncVisitor) isWanted(n ast.Node) bool {
    switch n.(type) {
        case *ast.File:
            return true
        case *ast.GenDecl:
            return true
        case *ast.FuncDecl:   
            return true
        case *ast.Comment:
            return true
        case *ast.ImportSpec:
            return true
        case *ast.ValueSpec:
            return true
        case *ast.TypeSpec:
            return true
        default:
            return false
    }
}

func (v *FuncVisitor) Visit(node ast.Node) ast.Visitor {
    if node == nil {
        p := asNode(v.stack.Pop())
        v.gotoOffset((*p).GetFooter(1))
    } else {
        if v.isWanted(node) {
            s := reflect.TypeOf(node).String()
            nodename := s[5:len(s)]
            n := v.createNode(node, nodename)
            v.stack.Push(n)
        } else {
            // Duplicating the top of the stack allows growing the DFS to stack to the right level and 
            // popping will correctly remove these duplicates. The effect on generated tree is that we
            // can skip levels:
            // This example of nesting: A->U1->U2->B (U* represent unwanted nodes) will be recorded
            // in the stack as:  A->A->A->B, therefore generating A->B
            v.stack.Push(v.stack.top.value) 
        }
    }
    return v
}

func parseErrors(filename string, v *FuncVisitor, errors scanner.ErrorList) {
    v.root = &File {
      id: newid(),
      Type: "file",
      Name: filename,
      ParsingErrorsDetected: true,
      FooterSpan: []int{ -1, 0 },
    }
    for _, error := range errors {
      loc := error.Pos;
      e := &ParsingError{ Location: []int{ loc.Line, loc.Column }, Message: error.Msg }
      v.root.ParsingError = append(v.root.ParsingError, e)
    }
}

func parseGoFile(filename string) (*FuncVisitor) {
    // Read as bytes buffer
    buffer, err1 := ioutil.ReadFile(filename)
    if err1 != nil {
        panic(err1)
    }
    // Parse file (get first file from returned fileset)
    fset := token.NewFileSet()
    file, errList := parser.ParseFile(fset, filename, nil, 0)
    var firstFile *token.File
    fset.Iterate(func(f *token.File) bool {
      firstFile = f
      return false
    })

    // Create stack for dfs
    stack := &Stack{}
    var v = &FuncVisitor{ stack: stack, file: firstFile, buffer: buffer }

    // Parse errors
    if errList != nil {
      parseErrors(filename, v, errList.(scanner.ErrorList))
    
    } else {
      // Do walk!
      ast.Walk(v, file)

      // Adjust FooterSpan with spare whitespaces at the end of the file
      endOfContent := v.adjustToEOL(v.root.FooterSpan[0]);
      endOfFile := len(buffer) - 1
      if endOfFile > endOfContent {
        v.root.FooterSpan = []int{ endOfContent + 1, endOfFile }
      } else {
        v.root.FooterSpan = []int{ 0, -1 }
      }
      n := len(v.root.Children)
      if n > 0 {
        last := v.root.Children[n - 1]
        lastChild := asNode(last)
        switch c := (*lastChild).(type) {
          case *Container:
            c.FooterSpan[1] = v.adjustToEOL(c.FooterSpan[1])
          case *Terminal:
            c.Span[1] = v.adjustToEOL(c.Span[1])
        }
      }
      // Dealing with transformation is a bit complicated in the first pass...
      // we transform the offsets to runes in a second pass of the output tree
      offsetTransformation(v.root, v.buffer)
    }    
    return v
}


//TODO optimize this!
func toRunes(offsets []int, buffer []byte) []int {
  return []int { bytesToRunes(offsets[0], buffer), bytesToRunes(offsets[1], buffer) }
}
func bytesToRunes(byteoffset int, buffer []byte) int {
  return utf8.RuneCount(buffer[0: byteoffset + 1]) - 1
}

func offsetTransformation(node interface{}, buffer []byte) {
  switch n := node.(type) {
    case *File:
      for _, c := range n.Children {
        offsetTransformation(c, buffer)
      }
      n.FooterSpan = toRunes(n.FooterSpan, buffer)
    
    case *Container:
      n.HeaderSpan = toRunes(n.HeaderSpan, buffer)
      for _, c := range n.Children {
        offsetTransformation(c, buffer)
      }
      n.FooterSpan = toRunes(n.FooterSpan, buffer)
    
    case *Terminal:
      n.Span = toRunes(n.Span, buffer)
      
    default:
      panic("Unknown Node! ")
  }
}

func WriteYAML(data []byte, outputFile string){
    f, err := os.Create(outputFile)
    if err != nil {
        fmt.Println("KO")
    }
    defer f.Close()
    _, err = f.Write(data)
    if err != nil {
        fmt.Println("KO")
    } else {
        fmt.Println("OK")
    }
}

func main() {
    syntax := "syntax error, please use:  goparser shell <flag file>"
    if len(os.Args) < 2 {
        fmt.Println(syntax)
        os.Exit(-1)
    }
    // there are two arguments to consider:
    // 1) "shell" saying you must run in "shell mode" 
    //    - don't exit basically and wait for commands
    // 2) A "flag file" - initialization is over
    shell := os.Args[1]
    if shell != "shell" {
        fmt.Println(syntax)
        os.Exit(-1)
    }
    flagFile := os.Args[2]

    // Write flag file immediately
    f, err := os.Create(flagFile)
    if err != nil {
        fmt.Println("KO")
        os.Exit(0)
    }
    _, err = f.WriteString("READY")
    f.Close()

    // Read STDIN
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        // read the file to parse first
        fileToParse := scanner.Text()
        if fileToParse == "end" {
            os.Exit(0)
        }

        // then where to put the resulting tree
        scanner.Scan()
        outputFile := scanner.Text()

        // Parse and marshall to file
        v := parseGoFile(fileToParse)
        
        d, err := yaml.Marshal(v.root)
        if err != nil {
            fmt.Println("KO")
        } else {
            // Write YAML file
            WriteYAML(d, outputFile)
        }
    }
}
