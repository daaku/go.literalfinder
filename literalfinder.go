// Package literalfinder helps find literals in source creating instances of a
// specified struct.
//
// This library uses the go/types library available in Go 1.1. The basic
// process is to create a Finder instance, add files from a package and find
// the contained literals. The API disallows non-literal instances. It
// currently also disallows positional initialization, though this may be
// allowed in the future. Additionally the struct may currently only contain
// basic literal values (bool, string, int64).
package literalfinder

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strconv"
)

var (
	errMustUseKeyValueSyntax = errors.New("must use key value value syntax")
)

// A Finder instance allows accumulating source files in a given package, thru
// which one can eventually find literal instances.
type Finder struct {
	thing   string
	fileSet *token.FileSet
	files   []*ast.File
}

// Create a new Finder instance. "thing" should be a fully qualified reference,
// like "github.com/go.pkgrsrc/pkgrsrc.Config". In that example,
// "github.com/go.pkgrsrc/pkgrsrc" is the package import path, and "Config" is
// the name of the struct type.
func NewFinder(thing string) *Finder {
	return &Finder{
		fileSet: token.NewFileSet(),
		thing:   thing,
	}
}

// Add a file from the package being processed. The src argument can be a
// string, []byte or io.Reader. If src == nil, the source will be read from the
// specified filename.
func (f *Finder) Add(filename string, src interface{}) error {
	astf, err := parser.ParseFile(f.fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	f.files = append(f.files, astf)
	return nil
}

// Find the instances and populate into. into should be an slice of the struct
// type that is being found.
func (f *Finder) Find(into interface{}) error {
	var instances []map[string]interface{}
	var retErr error
	exprFn := func(x ast.Expr, typ types.Type, val interface{}) {
		t, ok := typ.(*types.NamedType)
		if !ok {
			return
		}
		if t.String() != f.thing {
			return
		}
		l, ok := x.(*ast.CompositeLit)
		if !ok {
			return
		}
		fields, err := keyValueExprMap(l.Elts)
		if err != nil {
			retErr = err
			return
		}
		instances = append(instances, fields)
	}
	context := types.Context{
		Error: func(err error) {
			retErr = err
		},
		Expr: exprFn,
	}
	_, err := context.Check(f.fileSet, f.files)
	if err != nil {
		return err
	}
	if retErr != nil {
		return retErr
	}

	encoded, err := json.Marshal(instances)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(encoded, into); err != nil {
		return err
	}
	return nil
}

func keyValueExprMap(elts []ast.Expr) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, e := range elts {
		kv, ok := e.(*ast.KeyValueExpr)
		if !ok {
			return nil, errMustUseKeyValueSyntax
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			return nil, fmt.Errorf("unknown key type: %T", kv.Key)
		}
		val, err := literalValue(kv.Value)
		if err != nil {
			return nil, err
		}
		result[key.Name] = val
	}
	return result, nil
}

func literalValue(v ast.Expr) (interface{}, error) {
	switch i := v.(type) {
	case *ast.Ident:
		switch i.Name {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}
	case *ast.BasicLit:
		switch i.Kind {
		case token.STRING:
			l := len(i.Value)
			return i.Value[1 : l-1], nil
		case token.INT:
			return strconv.ParseInt(i.Value, 0, 64)
		case token.FLOAT:
			return strconv.ParseFloat(i.Value, 64)
		}
	}
	return nil, fmt.Errorf("unknown value type: %T", v)
}
