package main

import (
  "fmt"
  "testing"
  "bytes"
  "strconv"
  "os"
  "strings"
  "path/filepath"
)

const dump = false

func indent(level int) {
  for i := 0; i < level; i++ {
    fmt.Print("  ")
  }
}

func rebuild(node interface{}, v *FuncVisitor, runes []rune, output *bytes.Buffer, level int) bool {
  indent(level)
  switch n := node.(type) {
    case *File:
      //fmt.Println(node.GetName())
      for _, c := range n.Children {
        rebuild(c, v, runes, output, level + 1)
      }
      s := n.FooterSpan
      if s[0] > 0 {
        str := string(runes[s[0]:s[1]+1])
        fmt.Println(strconv.QuoteToASCII(str))
        output.WriteString(str)
      }
    case *Container:
      s := n.HeaderSpan
      str := string(runes[s[0]:s[1]+1])
      fmt.Println(strconv.QuoteToASCII(str))
      output.WriteString(str)
      for _, c := range n.Children {
        rebuild(c, v, runes, output, level + 1)
      }
      s = n.FooterSpan
      if s[0] > 0 {
        str = string(runes[s[0]:s[1]+1])
        fmt.Println(strconv.QuoteToASCII(str))
        output.WriteString(str)
      }
    case *Terminal:   
      s := n.Span
      str := string(runes[s[0]:s[1]+1])
      fmt.Println(s, strconv.QuoteToASCII(str))
      output.WriteString(str)
    default:
      panic("Unkown Node!")
  }
  return true
}

func doFile(filename string, t *testing.T) {
  v := parseGoFile(filename)
  var buffer *bytes.Buffer = &bytes.Buffer{}
  runes := []rune(string(v.buffer))
  
  rebuild(v.root, v, runes, buffer, 0)
  exp := v.buffer
  res := buffer.Bytes()

  // Dump if requested
  if dump {
    f, err := os.Create("testdata/dump.txt")
    defer f.Close()
    if err != nil {
        t.Error("Cannot read file", filename)
        return
    }
    f.Write(res)
  }
  
  // Check buffer
  if len(exp) != len(res) {
    t.Error("Filename: ", filename, ", expected length: ", len(exp), ", got ", len(res) )
    return
  }
  n := len(exp)
  for i := 0; i < n; i++ {
    if res[i] != exp[i] {
      t.Error("Filename: ", filename, ", expected ", exp[i], ", got ", res[i], " at offset: ", i)
      return
    }
  }
}

func TestWinBOM(t *testing.T) {
  doFile("testdata/example-win-bom.go", t)
}

func TestAll(t *testing.T) {
  filepath.Walk("testdata", func(path string, f os.FileInfo, err error) error {
    if strings.HasSuffix(path, ".go") {
      doFile(path, t)
    } 
    return nil
  });
}
