package lint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "explicittimelocation",
	Doc:  "explicittimelocation checks if time.Time has explicit location",
	Run:  run,
	// See https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/inspect
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func isTimeUnix(n ast.Node) (name string, ok bool) {
	if callExpr, ok := n.(*ast.CallExpr); ok {
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selectorExpr.X.(*ast.Ident); ok {
				if ident.Name == "time" {
					name := selectorExpr.Sel.Name
					if name == "Unix" || name == "UnixMicro" || name == "UnixMilli" {
						return name, true
					}
				}
			}
		}
	}
	return
}

func isUTC(n ast.Node) bool {
	if selectorExpr, ok := n.(*ast.SelectorExpr); ok {
		if selectorExpr.Sel.Name == "UTC" {
			return true
		}
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	traverse := func(n ast.Node, push bool, stack []ast.Node) (proceed bool) {
		if push == false {
			shouldCheck := false
			idx := -1
			funcName := ""

			for i := len(stack) - 1; i >= 0; i -= 1 {
				node := stack[i]
				var ok bool
				funcName, ok = isTimeUnix(node)
				if ok {
					shouldCheck = true
					idx = i
					break
				}
			}
			if shouldCheck {
				node := stack[idx]
				parentNodeIdx := idx - 1
				if parentNodeIdx < 0 {
					pass.Reportf(node.Pos(), "time.%v() is not immediately followed by .UTC()", funcName)
				} else {
					parentNode := stack[parentNodeIdx]
					if !isUTC(parentNode) {
						pass.Reportf(node.Pos(), "time.%v() is not immediately followed by .UTC()", funcName)
					}
				}
			}
		}
		return true
	}
	inspect.WithStack([]ast.Node{(*ast.CallExpr)(nil)}, traverse)
	return nil, nil
}
