package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var InterfaceMustBePtr = &analysis.Analyzer{
	Name:     "interfacemustbeptr",
	Doc:      "Checks that calls that take an interface{} are passed a pointer",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var interfaceMustBePtrTargets = map[call][]int{
	{pkg: "encoding/json", fxn: "Unmarshal"}: {1},
}

type call struct {
	pkg, fxn string
}

func run(pass *analysis.Pass) (interface{}, error) {

	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var imports map[string]string

	var currentScope *types.Scope

	isPointer := func(t types.Type) bool {
		return t.String()[0] == '*'
	}
	scope2Node := func(target *types.Scope) ast.Node {
		for node, scope := range pass.TypesInfo.Scopes {
			if scope == target {
				return node
			}
		}
		return nil
	}

	inspect := func(node ast.Node, push bool, stack []ast.Node) bool {
		if node == nil {
			return true
		}
		if push == false {
			return true
		}

		fmt.Println("=", node, push, reflect.TypeOf(node))
		scope, hasScope := pass.TypesInfo.Scopes[node]
		if hasScope {
			currentScope = scope
		}

		if file, ok := node.(*ast.File); ok {
			fmt.Println("🎸", file.Name, push)
			imports = make(map[string]string)
			return true
		}

		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			fmt.Println("___", funcDecl.Name.Name, "____")
			return true
		}

		if importSpec, ok := node.(*ast.ImportSpec); ok {
			value := importSpec.Path.Value
			value = value[1 : len(value)-1] // strip quotes
			var name string
			if importSpec.Name != nil {
				name = importSpec.Name.Name
			} else {
				name = filepath.Base(value)
			}
			imports[name] = value
			return true
		}

		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		// fmt.Println("__", callExpr, callExpr.Fun, reflect.TypeOf(callExpr.Fun))

		selector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		// fmt.Println("selector name", selector.Sel.Name)
		x := selector.X.(*ast.Ident)
		// TODO handle methods x.Obj :)
		// fmt.Println("x", x, reflect.TypeOf(x))
		// fmt.Println("x.Name", x.Name, reflect.TypeOf(x.Name))
		pkg, ok := imports[x.Name]
		if !ok {
			err := fmt.Errorf("couldn't resolve identifier %s to import", x.Name)
			panic(err)
		}
		c := call{pkg: pkg, fxn: selector.Sel.Name}
		argNumbers, ok := interfaceMustBePtrTargets[c]
		if ok {
			// fmt.Printf("matched! %s.%s\n", c.pkg, c.fxn)
		} else {
			return true
		}

		fmt.Println("(")
		for _, argNumber := range argNumbers {
			arg := callExpr.Args[argNumber]
			fmt.Println("arg", argNumber, arg, reflect.TypeOf(arg))
			t := pass.TypesInfo.TypeOf(arg)
			isPtr := isPointer(t)
			fmt.Println("arg isPtr?", isPtr)
			if ident, ok := arg.(*ast.Ident); ok {
				o := pass.TypesInfo.ObjectOf(ident)
				fmt.Println("ident is object of", o)
				fmt.Println("   ", o.Parent())
				fmt.Println("decl", ident.Obj.Decl, reflect.TypeOf(ident.Obj.Decl))
				if decl, ok := ident.Obj.Decl.(*ast.Field); ok {
					fmt.Println("pos", decl.Pos())
					// pass.TypesInfo.ObjectOf(decl)

					from := scope2Node(o.Parent())
					if from != nil {
						fmt.Println("from   ", (from), reflect.TypeOf(from))
						// from.(*ast.FuncType)
						if fxn, ok := from.(*ast.FuncType); ok {
							for _, item := range fxn.Params.List {
								fmt.Println("item", item, reflect.TypeOf(item))
								for i, name := range item.Names {
									fmt.Println("name", name, reflect.TypeOf(name))
									if name.Name == decl.Names[0].Name {
										fmt.Println("name is taint", i, name)
									}
								}
							}
							fmt.Println("params", fxn.Params.NumFields())
						}
					}
				}
				continue
				fmt.Println("ident name", ident.Name)
				fmt.Println("ident obj", ident.Obj)
				fmt.Println("ident obj kind", ident.Obj.Kind, reflect.TypeOf(ident.Obj.Kind))
				fmt.Println("ident obj decl", ident.Obj.Decl, reflect.TypeOf(ident.Obj.Decl))
				if currentScope == nil {
					panic("nil scope")
				}
				_, obj := currentScope.LookupParent(ident.Name, ident.NamePos)
				if obj == nil {
					fmt.Println("LookupParent is nil") // error
				} else {
					isPtr := isPointer(obj.Type())
					fmt.Println("obj isPtr?", isPtr)
				}
			} else if callExpr, ok := arg.(*ast.CallExpr); ok {
				fmt.Println("callExpr", callExpr)
			}
			fmt.Println(",")
		}
		fmt.Println(")")

		return true
	}

	inspector.WithStack(nil, inspect)

	return nil, nil
}
