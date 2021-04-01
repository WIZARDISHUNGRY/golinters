package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var MarshalPlan = &analysis.Analyzer{
	Name:     "marshalplan",
	Doc:      "Checks that calls that take an interface{} are passed a pointer",
	Run:      (&marshalPlan{}).Run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type marshalPlan struct {
	currentScope *types.Scope
	pass         *analysis.Pass
}

func (mp *marshalPlan) Run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	mp.pass = pass

	inspect := func(node ast.Node, push bool, stack []ast.Node) bool {
		if node == nil {
			return true
		}
		if push == false {
			return true
		}

		scope, hasScope := pass.TypesInfo.Scopes[node]
		if hasScope {
			mp.currentScope = scope
		}

		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			ret := funcDecl.Name.String() == "Fail_UnmarshalMap_Indirect2" || funcDecl.Name.String() == "indirect" || funcDecl.Name.String() == "Pass_UnmarshalMap_Indirect2"
			if ret {
				fmt.Println("___", funcDecl.Name.Name, "____")
			}
			return ret
			return true
		}

		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		fmt.Println("call__", callExpr, callExpr.Fun, reflect.TypeOf(callExpr.Fun), mp.currentScope)

		fun := callExpr.Fun
		flagFunction := false
		switch fun.(type) {
		case *ast.ArrayType:
		case *ast.Ident:
			f := mp.ident(node, fun.(*ast.Ident))
			fmt.Println(f)
			flagFunction = true
		case *ast.SelectorExpr:
			flagFunction = mp.selectorExpr(node, fun.(*ast.SelectorExpr))
		case *ast.FuncLit:
		default:
			panic(fun)
		}

		if flagFunction {
			fmt.Println("danger scope", mp.currentScope)
			for i, arg := range callExpr.Args {
				t := pass.TypesInfo.TypeOf(arg)
				fmt.Println("arg", i, arg, reflect.TypeOf(arg))
				fmt.Println("   ", i, t, reflect.TypeOf(t))
			}

		}

		return true
	}

	inspector.WithStack(nil, inspect)

	return nil, nil
}
func (mp *marshalPlan) arrayType(callExpr *ast.ArrayType) {
}
func (mp *marshalPlan) ident(node ast.Node, callExpr *ast.Ident) *types.Func {
	_, o := mp.currentScope.LookupParent(callExpr.Name, node.Pos())

	if f, ok := o.(*types.Func); ok {
		fmt.Println("🌹", o, reflect.TypeOf(o), f.Type(), reflect.TypeOf(f.Type()))
		if sig, ok := f.Type().(*types.Signature); ok {

			for i := 0; i < sig.Params().Len(); i++ {
				v := sig.Params().At(i)
				tv := mp.pass.TypesInfo.Types[callExpr]
				fmt.Println("X", v, "🎨", tv, "🎨")

				// mp.currentScope.Lookup()
			}
		}
		return f
	}
	return nil
}
func (mp *marshalPlan) selectorExpr(node ast.Node, callExpr *ast.SelectorExpr) bool {
	fmt.Println("🚴🏿", callExpr.Sel)                           // marshal
	fmt.Println("🚴🏿", callExpr.X, reflect.TypeOf(callExpr.X)) // json
	pkgName := callExpr.X.(*ast.Ident).Name
	_, o := mp.currentScope.LookupParent(pkgName, node.Pos())
	fmt.Println("🚴🏿", callExpr.Sel.Obj, o, reflect.TypeOf(o))
	return true
}

func (mp *marshalPlan) funcLit(callExpr *ast.FuncLit) {}
