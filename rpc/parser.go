package rpc

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gomicro-tools/common"
)

type Var struct {
	IsSlice bool
	Type    string
	Name    string
}

type Method struct {
	Name       string
	Parameters []*Var
	Returns    []*Var
}

type InterFace struct {
	Name    string
	Methods []*Method
}

type Struct struct {
	Name    string
	Members []*Var
}

// TODO: matrix
func parseType(typeExpr ast.Expr) (string, bool) {
	switch expr := typeExpr.(type) {
	case *ast.Ident: // int, string, error, ...
		return expr.Name, false
	case *ast.StarExpr: // *xxx
		typeName, _ := parseType(expr.X)
		return typeName, false
	case *ast.ArrayType: // []xxx
		typeName, _ := parseType(expr.Elt)
		return typeName, true
	case *ast.SelectorExpr: // pkg.Foo, 返回 Foo
		return expr.Sel.Name, false
	default:
		panic(fmt.Sprintf("unexcepted type %[1]T: %[1]v (%#[1]v)", expr))
	}
}

func parseFieldList(fieldList []*ast.Field) []*Var {
	var vars []*Var

	for _, field := range fieldList {
		typeName, isSlice := parseType(field.Type)

		for _, nameIdent := range field.Names {
			vars = append(vars, &Var{isSlice, typeName, nameIdent.Name})
		}
	}

	return vars
}

func genAstFile(sourceCode string) *ast.File {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "", sourceCode, 0)
	common.Check(err)

	return f
}

func ParseInterface(sourceCode string) *InterFace {
	f := genAstFile(sourceCode)

	for _, v := range f.Scope.Objects {
		typeSpec, ok := v.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		interfaceName := typeSpec.Name.Name

		methodFieldList := interfaceType.Methods.List
		parsedMethods := make([]*Method, len(methodFieldList))

		for i, field := range methodFieldList {
			funcType, ok := field.Type.(*ast.FuncType)
			if !ok {
				continue
			}

			methodName := field.Names[0].Name

			parsedMethods[i] = &Method{
				methodName,
				parseFieldList(funcType.Params.List),
				parseFieldList(funcType.Results.List),
			}
		}

		return &InterFace{interfaceName, parsedMethods}
	}

	return nil
}

func parseStructs(sourceCode string) []*Struct {
	f := genAstFile(sourceCode)

	var parsedStructs []*Struct
	for _, v := range f.Scope.Objects {
		typeSpec, ok := v.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}

		structName := typeSpec.Name.Name

		structFieldList := structType.Fields.List
		members := parseFieldList(structFieldList)

		parsedStructs = append(parsedStructs, &Struct{structName, members})
	}
	return parsedStructs
}

func parseStructsFromCodes(sourceCode []string) []*Struct {
	var parsedStructs []*Struct

	for _, src := range sourceCode {
		parsedStructs = append(parsedStructs, parseStructs(src)...)
	}

	return parsedStructs
}
