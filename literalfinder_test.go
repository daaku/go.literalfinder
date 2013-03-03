package literalfinder_test

import (
	"github.com/daaku/go.literalfinder"
	"testing"
)

func TestNoLiterals(t *testing.T) {
	t.Parallel()
	const source = `
  package foo
  type Foo struct {
    Bar bool
  }
  `
	i, err := literalfinder.Find("Foo", "foo.go", source)
	if err != nil {
		t.Fatal(err)
	}
	if len(i) != 0 {
		t.Fatal("was expecting 0 instance")
	}
}

func TestSingleString(t *testing.T) {
	t.Parallel()
	const source = `
  package foo
  type Foo struct {
    Bar string
  }
  var f = &Foo{Bar: "one"}
  `
	i, err := literalfinder.Find("Foo", "foo.go", source)
	if err != nil {
		t.Fatal(err)
	}
	if len(i) != 1 {
		t.Fatal("was expecting 1 instance")
	}
	if v := i[0].Fields["Bar"]; v != "one" {
		t.Fatalf("was expecting one got %s", v)
	}
}

func TestBools(t *testing.T) {
	t.Parallel()
	const source = `
  package foo
  type Foo struct {
    Bar bool
  }
  var f = &Foo{Bar: true}
  var g = &Foo{Bar: false}
  `
	i, err := literalfinder.Find("Foo", "foo.go", source)
	if err != nil {
		t.Fatal(err)
	}
	if len(i) != 2 {
		t.Fatal("was expecting 2 instance")
	}
	if v := i[0].Fields["Bar"]; v != true {
		t.Fatalf("was expecting true got %s", v)
	}
	if v := i[1].Fields["Bar"]; v != false {
		t.Fatalf("was expecting false got %s", v)
	}
}
