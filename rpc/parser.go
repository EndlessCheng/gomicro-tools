package rpc

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type Var struct {
	IsSlice bool
	Type    string
	Name    string
}

type Method struct {
	MethodName string
	Parameters []*Var
	Returns    []*Var
}

type InterFace struct {
	InterFaceName string
	Methods       []*Method
}

type Struct struct {
	StructName string
	Members    []*Var
}

func parseType(typeExpr ast.Expr) (string, bool) {
	ident, ok := typeExpr.(*ast.Ident)
	if ok {
		return ident.Name, false
	}

	arrayType, ok := typeExpr.(*ast.ArrayType)
	if ok {
		typeName, _ := parseType(arrayType.Elt)
		// TODO: matrix
		return typeName, true
	}

	//panic(typeExpr)
	return "", false
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

func parseInterface(sourceCode string) *InterFace {
	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, "", sourceCode, 0)
	check(err)

	for _, v := range f.Scope.Objects {
		typeSpec, ok := v.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

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

		return &InterFace{v.Name, parsedMethods}
	}

	return nil
}

func parseStruct(srcFilePath string) *Struct {
	return nil
}
