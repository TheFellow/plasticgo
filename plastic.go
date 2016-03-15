package main

import "reflect"

type Location struct {
    Start []int  `yaml:"start,flow"`
    End []int    `yaml:"end,flow"`
}

type YamlNode interface {
    GetId() uint
    GetName() string
    AddChild(child interface{})
    GetFooter(idx int) int
    GetParent() interface{}
}

type File struct {
    id uint
    Type string                  `yaml:"type"` // Always "file"
    Name string                  `yaml:"name"` // path of the file
    LocationSpan Location        `yaml:"locationSpan,omitempty,flow"` // row and column where the file starts and ends (optional)
    FooterSpan []int             `yaml:"footerSpan,flow"` // start and end char where the file starts and ends
    ParsingErrorsDetected bool   `yaml:"parsingErrorsDetected"` // whether or not the file contains parsing errors
    Children []interface{}       `yaml:"children,omitempty"` // set of containers and/or terminal nodes inside the file. If there aren't any, this field shouldn't be specified.
    ParsingError []*ParsingError `yaml:"parsingError,omitempty"` // set of parsing errors (optional, see description below)
}

type Container struct {
    id uint
    parent interface{}
    Type string               `yaml:"type"` // relevant, generic name of the container in the current programming language
    Name string               `yaml:"name"` // actual name of the container
    LocationSpan *Location    `yaml:"locationSpan,omitempty,flow"` // row and column where the container starts and ends (optional)
    HeaderSpan []int          `yaml:"headerSpan,flow"` // start and end chars where the header of the container starts and ends
    FooterSpan []int          `yaml:"footerSpan,flow"` // start and end chars where the footer of the container starts and ends. This field should be set to [0, -1] if unexisting
    Children []interface{}    `yaml:"children,omitempty"` // set of containers and/or terminal nodes present inside the current container. If there aren’t any, this field shouldn’t be specified.
}

type Terminal struct {
    id uint
    parent interface{}
    Type string               `yaml:"type"` // relevant, generic name of the node in the current programming language
    Name string               `yaml:"name"` // actual name of the node
    LocationSpan *Location    `yaml:"locationSpan,omitempty,flow"` // row and column where the node starts and ends (optional)
    Span []int                `yaml:"span,flow"` // start and end char where the node starts and ends   
}

type ParsingError struct {
    Location []int            `yaml:"location,flow"`
    Message string            `yaml:"message"`
}

func addChild(parent interface{}, child interface{}) {
  switch n := parent.(type) {
    case *File:
      n.Children = append(n.Children, child)
    case *Container:
      n.Children = append(n.Children, child)
    default:
      s := reflect.TypeOf(n).String()
      panic("unknown parent: " + s)
  }
  switch n := child.(type) {
    case *Container:
      n.parent = parent
    case *Terminal:
      n.parent = parent
    default:
      s := reflect.TypeOf(n).String()
      panic("unknown child: " + s)
  }
}

func (f *File) GetId() uint { return f.id }
func (f *File) GetName() string { return f.Name }
func (f *File) GetParent() interface{} { return nil }
func (f *File) AddChild(child interface{}) { addChild(f, child) }
func (f *File) GetFooter(idx int) int { return f.FooterSpan[idx] }

func (c *Container) GetId() uint { return c.id }
func (c *Container) GetName() string { return c.Name }
func (c *Container) GetParent() interface{} { return c.parent }
func (c *Container) AddChild(child interface{}) { addChild(c, child) }
func (c *Container) GetFooter(idx int) int { return c.FooterSpan[idx] }

func (t *Terminal) GetId() uint { return t.id }
func (t *Terminal) GetName() string { return t.Name }
func (t *Terminal) GetParent() interface{} { return t.parent }
func (t *Terminal) AddChild(child interface{}) { }
func (t *Terminal) GetFooter(idx int) int { return t.Span[idx] }
