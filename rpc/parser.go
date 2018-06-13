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
	// int, string, ...
	ident, ok := typeExpr.(*ast.Ident)
	if ok {
		return ident.Name, false
	}

	// *xxx
	starExpr, ok := typeExpr.(*ast.StarExpr)
	if ok {
		typeName, _ := parseType(starExpr.X)
		return typeName, false
	}

	// []xxx
	arrayType, ok := typeExpr.(*ast.ArrayType)
	if ok {
		typeName, _ := parseType(arrayType.Elt)
		return typeName, true
	}

	// pkg.Foo
	selectorExpr, ok := typeExpr.(*ast.SelectorExpr)
	if ok {
		return selectorExpr.Sel.Name, false
	}

	// ?
	panic(typeExpr)
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
	check(err)

	return f
}

func parseInterface(sourceCode string) *InterFace {
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

		parsedStructs = append(parsedStructs, &Struct{structName,members})
	}
	return parsedStructs
}
