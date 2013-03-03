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
	if err := literalfinder.Find(&foos, "Foo", "foo.go", source); err != nil {
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
	if err := literalfinder.Find(&foos, "Foo", "foo.go", source); err != nil {
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

func TestNoLiterals(t *testing.T) {
	t.Parallel()
	var foos []struct{ Bar bool }
	const source = `
  package foo
  type Foo struct { Bar bool }
  `
	if err := literalfinder.Find(&foos, "Foo", "foo.go", source); err != nil {
		t.Fatal(err)
	}
	if len(foos) != 0 {
		t.Fatal("was expecting 0 instance")
	}
}
