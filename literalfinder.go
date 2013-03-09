// Package literalfinder helps find literals in source creating instances of a
// specified struct.
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

type Finder struct {
	thing   string
	fileSet *token.FileSet
	files   []*ast.File
}

func NewFinder(thing string) *Finder {
	return &Finder{
		fileSet: token.NewFileSet(),
		thing:   thing,
	}
}

func (f *Finder) Add(filename string, src interface{}) error {
	astf, err := parser.ParseFile(f.fileSet, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	f.files = append(f.files, astf)
	return nil
}

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
		default:
			return i.Value, nil
		}
	}
	return nil, fmt.Errorf("unknown value type: %T", v)
}
