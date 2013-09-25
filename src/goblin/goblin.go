package goblin

import (
  "testing"
  "fmt"
)

var parentDescribe *D

type Runnable interface {
  Run() (bool)
}

type It struct {
  name string
  h func(*T)
  t *T
  parent *D
}

func (it It) Run() (bool) {
  //TODO: should handle errors for beforeEach
  it.parent.runBeforeEach()

  fmt.Println(it.name)
  it.h(it.t)

  return !it.t.Failed()
}

type D struct {
  name string
  parent *D
  children []Runnable
  beforeEach []func()
}

func Describe(name string, h func(*D)) {
  parentDescribe = &D{name: name}
  h(parentDescribe)
}

func (d *D) Describe(name string, h func(*D)) {
  describe := &D{name: name, parent: d}
  d.addChild(Runnable(describe))
  h(describe)
}

func (d *D) runBeforeEach() {
  if d.parent != nil {
    d.parent.runBeforeEach()
  }

  for _, b := range d.beforeEach {
    b()
  }
}

func (d *D) BeforeEach(h func()) {
  d.beforeEach = append(d.beforeEach, h)
}

func (d *D) addChild(r Runnable) {
  d.children = append(d.children, r)
}

func (d *D) It(name string, h func(t *T)) {
  it := It{name: name, h: h, t: &T{}, parent: d}
  d.addChild(Runnable(it))
}

func (d D) Run() (bool) {
  succeed := true
  for _, r := range d.children {
    if !r.Run() {
      succeed = false
    }
  }
  //TODO: run afterEach
  return succeed
}


func Goblin(t *testing.T) {
  succeed := parentDescribe.Run()

  if !succeed {
    t.Fail()
  }
}
