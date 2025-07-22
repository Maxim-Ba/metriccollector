package staticlint

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// OsExitAnalyzer checks for direct calls to os.Exit in the main function of main packages.
//
// This analyzer detects and reports usage of os.Exit() in main functions, which is
// generally discouraged as it can prevent proper cleanup of resources and bypass
// deferred functions. Instead, it's recommended to return from main or use proper
// error handling patterns.
//
// Usage: multichecker [-flag] [package]
//
// Example of a finding:
//
//	func main() {
//	    os.Exit(1) // this will be reported
//	}
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for using os.Exit in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		// Проверяем, что это пакет main
		if pass.Pkg.Name() != "main" {
			continue
		}

		// Ищем функцию main
		ast.Inspect(file, func(n ast.Node) bool {
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				return true // продолжаем поиск
			}

			// Проверяем тело функции main на вызовы os.Exit
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				// Ищем выражения-вызовы функций
				expr, ok := n.(*ast.ExprStmt)
				if !ok {
					return true
				}

				// Проверяем, является ли выражение вызовом функции
				call, ok := expr.X.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Проверяем, является ли это вызовом os.Exit
				if isOsExitCall(pass, call) {
					pass.Reportf(call.Pos(), "direct call to os.Exit in main function of main package")
				}

				return true
			})

			return false // уже нашли main, дальше не ищем
		})
	}

	return nil, nil
}

// isOsExitCall  возвращает true, если это вызов os.Exit.
func isOsExitCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false // это не селектор (не вызов вида pkg.Func)
	}

	// Проверяем, что это вызов из пакета os
	if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "os" {
		if sel.Sel.Name == "Exit" {
			if fn, ok := pass.TypesInfo.ObjectOf(sel.Sel).(*types.Func); ok {
				return fn.Pkg().Path() == "os" && fn.Name() == "Exit"
			}
		}
	}

	return false
}
