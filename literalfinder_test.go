package literalfinder_test

import (
	"github.com/daaku/go.literalfinder"
	"testing"
)

func TestSingleString(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar string }
	const source = `
  package foo
  type Foo struct { Bar string }
  var f = &Foo{Bar: "one"}
  `
	f := literalfinder.NewFinder("Foo")
	f.Add("foo.go", source)
	if err := f.Find(&foos); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 1 {
		t.Fatal("was expecting 1 instance")
	}
	if v := foos[0].Bar; v != "one" {
		t.Fatalf("was expecting one got %s", v)
	}
}

func TestBools(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar bool }
	const source = `
  package foo
  type Foo struct { Bar bool }
  var f = &Foo{Bar: true}
  var g = &Foo{Bar: false}
  `
	f := literalfinder.NewFinder("Foo")
	f.Add("foo.go", source)
	if err := f.Find(&foos); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 2 {
		t.Fatal("was expecting 2 instance")
	}
	if v := foos[0].Bar; v != true {
		t.Fatalf("was expecting true got %s", v)
	}
	if v := foos[1].Bar; v != false {
		t.Fatalf("was expecting false got %s", v)
	}
}

func TestInt(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar int }
	const source = `
  package foo
  type Foo struct { Bar int }
  var f = &Foo{Bar: 42}
  `
	f := literalfinder.NewFinder("Foo")
	f.Add("foo.go", source)
	if err := f.Find(&foos); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 1 {
		t.Fatal("was expecting 1 instance")
	}
	if v := foos[0].Bar; v != 42 {
		t.Fatalf("was expecting one got %s", v)
	}
}

func TestFloat(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar float64 }
	const source = `
  package foo
  type Foo struct { Bar float64 }
  var f = &Foo{Bar: 4.2}
  `
	f := literalfinder.NewFinder("Foo")
	f.Add("foo.go", source)
	if err := f.Find(&foos); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 1 {
		t.Fatal("was expecting 1 instance")
	}
	if v := foos[0].Bar; v != 4.2 {
		t.Fatalf("was expecting one got %s", v)
	}
}

func TestNoLiterals(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar bool }
	const source = `
  package foo
  type Foo struct { Bar bool }
  `
	f := literalfinder.NewFinder("Foo")
	f.Add("foo.go", source)
	if err := f.Find(&foos); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 0 {
		t.Fatal("was expecting 0 instance")
	}
}
