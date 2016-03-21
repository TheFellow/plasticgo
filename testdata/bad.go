package main

import   "reflect" 
import m "math"
import . "math"

import ( ==
  _ "bytes"  // notused but ignored
  `fmt`      // alternatives quotes
  "go/ast"   // used
)

func myfunction() {
  // An empty struct.
  var primera struct {}

  // A struct with 6 fields. 
  var segunda struct {
    x, y NonExistent
    u float32
    _ float32  // padding
    A *[]int
    F func()
  }

  // Funtion types
  var a func()
  var b func(x int) int
  var c func(a, _ int, z float32) bool
  var d func(a, b int, z float32) (bool)
  var e func(prefix string, values ...int)
  var f func(a, b int, z float64, opt ...interface{}) (success bool)
  var g func(int, int, float64) (float64, *[]int)
  var h func(n int) func(p *Locker)
  
  _=primera
  _=segunda

  _=a
  _=b
  _=c
  _=d
  _=e
  _=f
  _=g
  _=h

  // A simple File interface
  var iii interface {
    Read(b Buffer) bool
    Write(b Buffer) bool
    Close()
  }
  
  _=iii

}


type IntArray [16]int

type (
	Point struct{ x, y float64 }
	Polar Point
)

// Be careful with this! should a inline list of types be a terminal or a container
type ( Point2 struct{ x, y float64 }; Polar2 Point2; )


type TreeNode struct {
	left, right *TreeNode
	value string
}

type Block interface {
	BlockSize() int
	Encrypt(src, dst []byte)
	Decrypt(src, dst []byte)
}

type Locker interface {
  LockerMethod()
}

type Buffer struct {
  Content []byte
}

type ReadWriter interface {
	Read(b Buffer) bool
	Write(b Buffer) bool
}

type File interface {
	ReadWriter  // same as adding the methods of ReadWriter
	Locker      // same as adding the methods of Locker
	Close()
}

const OtraPi float64 = 3.14159265358979323846
const zero = 0.0         // untyped floating-point constant
const (
	size int64 = 1024
	eof        = -1  // untyped integer constant
)
const a, b, c = 3, 4, "foo"  // a = 3, b = 4, c = "foo", untyped integer and string constants
const u, v float32 = 0, 3    // u = 0.0, v = 3.0

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Partyday
	numberOfDays  // this constant is not exported
)

const ( // iota is reset to 0
	c0 = iota  // c0 == 0
	c1 = iota  // c1 == 1
	c2 = iota  // c2 == 2
)

const ( // iota is reset to 0
	aa = 1 << iota  // a == 1
	bb = 1 << iota  // b == 2
	cc = 3          // c == 3  (iota is not used but still incremented)
	dd = 1 << iota  // d == 8
)

const ( // iota is reset to 0
	uu         = iota * 42  // u == 0     (untyped integer constant)
	vv float64 = iota * 42  // v == 42.0  (float64 constant)
	ww         = iota * 42  // w == 84    (untyped integer constant)
)

const x = iota  // x == 0  (iota has been reset)
const y = iota  // y == 0  (iota has been reset)

// A Mutex is a data type with two methods, Lock and Unlock.
type Mutex struct         { /* Mutex fields */ }
func (m *Mutex) Lock()    { /* Lock implementation */ }
func (m *Mutex) Unlock()  { /* Unlock implementation */ }

// NewMutex has the same composition as Mutex but its method set is empty.
type NewMutex Mutex

// The method set of the base type of PtrMutex remains unchanged,
// but the method set of PtrMutex is empty.
type PtrMutex *Mutex

// The method set of *PrintableMutex contains the methods
// Lock and Unlock bound to its anonymous field Mutex.
type PrintableMutex struct {
	Mutex
}

// MyBlock is an interface type that has the same method set as Block.
type MyBlock Block

type TimeZone int

const (
	EST TimeZone = -(5 + iota)
	CST
	MST
	PST
)

func (tz TimeZone) String() string {
	return fmt.Sprintf("GMT%+dh", tz)
}

var i int
var U, V, W float64
var k = 0
var xx, yy float32 = -1, -2
var (
	ii       int
	uuu, vvv, sss = 2.0, 3.0, "bar"
)

func IndexRune(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return -1
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

// func flushICache(begin, end uintptr)  // implemented externally

func (p *Point) Length() float64 {
	return m.Sqrt(p.x * p.x + p.y * p.y)
}

func (p *Point) Scale(factor float64) {
	p.x *= factor
	p.y *= factor
}

func main() {
	type Point3D struct { x, y, z float64 }
	type Line struct { p, q Point3D }
	origin := Point3D{}                            // zero value for Point3D
	line := Line{origin, Point3D{y: -4, z: 12.3}}  // zero value for line.q.x
	var pointer *Point3D = &Point3D{y: 1000}

	buffer := [10]string{}             // len(buffer) == 10
	intSet := [6]int{1, 2, 3, 5}       // len(intSet) == 6
	days := [...]string{"Sat", "Sun"}  // len(days) == 2
	
	// list of prime numbers
	primes := []int{2, 3, 5, 7, 9, 2147483647}
	
	// vowels[ch] is true if ch is a vowel
	vowels := [128]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'y': true}
	
	// the array [10]float32{-1, 0, 0, 0, -0.1, -0.1, 0, 0, 0, -1}
	filter := [10]float32{-1, 4: -0.1, -0.1, 9: -1}
	
	// frequencies in Hz for equal-tempered scale (A4 = 440Hz)
	noteFrequency := map[string]float32{
		"C0": 16.35, "D0": 18.35, "E0": 20.60, "F0": 21.83,
		"G0": 24.50, "A0": 27.50, "B0": 30.87,
	}
  var v1  chan int 
	f := func(x, y int) int { return x + y }
	func(ch chan int) { ch <- 123 } (v1)
  
  // Use imported packages
  sss := reflect.TypeOf(f).String()
  fmt.Println(m.Sin(0.4))
  fmt.Println(Sin(0.5))
  var xxx ast.BasicLit
  
  _=line
  _=pointer
  _=buffer
	_=intSet
  _=days
  _=primes
  _=vowels
  _=filter
  _=noteFrequency
  _=sss
  _=xxx
  
    var state string
    var a [][]int
    var i int
    var j int
    n := 10
    m := 9

  OuterLoop:
		for i = 0; i < n; i++ {
			for j = 0; j < m; j++ {
				switch a[i][j] {
				case 1:
					state = "Error"
					break OuterLoop
				case 2:
					state = "Found"
					break OuterLoop
				}
        fmt.Println("state", state)
			}
      
      
	} 

}


