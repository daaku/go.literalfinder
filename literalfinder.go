// Package literalfinder helps find literals in source creating instances of a
// specified struct.
package literalfinder

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

var (
	errMustUseKeyValueSyntax = errors.New("must use key value value syntax")
)

type Instance struct {
	Fields map[string]interface{}
}

func Find(thing string, filename string, src interface{}) ([]Instance, error) {
	fset := token.NewFileSet()
	astf, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	var instances []Instance
	exprFn := func(x ast.Expr, typ types.Type, val interface{}) {
		t, ok := typ.(*types.NamedType)
		if !ok {
			return
		}
		if t.String() != thing {
			return
		}
		l, ok := x.(*ast.CompositeLit)
		if !ok {
			//panic("non literal instance")
			return
		}
		fields, err := keyValueExprMap(l.Elts)
		if err != nil {
			panic(err)
		}
		instances = append(instances, Instance{Fields: fields})
	}
	context := types.Context{
		Expr: exprFn,
	}
	_, err = context.Check(fset, []*ast.File{astf})
	if err != nil {
		return nil, err
	}
	return instances, nil
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
		default:
			return i.Value, nil
		}
	}
	return nil, fmt.Errorf("unknown value type: %T", v)
}
