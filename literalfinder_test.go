package literalfinder_test

import (
	"github.com/daaku/go.literalfinder"
	"testing"
)

func TestSimpleFind(t *testing.T) {
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
