package usecase

import (
	// "bytes"
	"fmt"
	"go/ast"
	"go/parser"
	// "go/printer"
	"go/token"
	"strings"

	"github.com/suzuki-shunsuke/go-set"

	"github.com/suzuki-shunsuke/gria/domain"
)

func isExposed(name string) bool {
	return len(name) != 0 && name[:1] == strings.ToUpper(name[:1])
}

// Gen generates test scaffold.
func Gen() error {
	// test package name

	// read the file
	// get exposed function and method names
	// read the test file
	// get exposed function names
	// get added functions
	// create code
	// append code to test file
	fset := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fset, "dummy", nil, parser.Mode(0))
	funcs := []domain.Func{}
	testFuncs := domain.Funcs{Names: map[string]map[string]domain.Func{}}
	fileNameSet := set.NewStrSet()
	// file name -> appended test code
	addedFuncs := map[string]string{}
	for pkgName, pkg := range pkgs {
		for fileName, f := range pkg.Files {
			fileNameSet.Add(fileName)
			ast.Inspect(f, func(n ast.Node) bool {
				if ident, ok := n.(*ast.FuncDecl); ok {
					fName := ident.Name.Name
					if !isExposed(fName) {
						return true
					}
					fnc := domain.Func{
						Name:        fName,
						FileName:    fileName,
						PackageName: pkgName,
					}
					if fnc.IsTest() {
						testFuncs.Add(fnc)
						return true
					}
					if fnc.IsTestFile() {
						return true
					}
					funcs = append(funcs, fnc)
				}
				return true
			})
		}
	}
	// A_test.go B B_test TestC
	for _, f := range funcs {
		if _, ok := testFuncs.Names[f.Name]; ok {
			continue
		}
		if m, ok := addedFuncs[f.FileName]; ok {
			m += f.TestCode()
			addedFuncs[f.FileName] = m
			continue
		}
		addedFuncs[f.FileName] = f.TestCode()
	}

	for fName, code := range addedFuncs {
		fmt.Println(fName)
		fmt.Println(code)
	}
	return nil
}
