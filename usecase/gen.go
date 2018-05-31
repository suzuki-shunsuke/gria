package usecase

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/suzuki-shunsuke/go-set"

	"github.com/suzuki-shunsuke/gria/domain"
)

func isExposed(name string) bool {
	return len(name) != 0 && name[:1] == strings.ToUpper(name[:1])
}

// Gen generates test scaffold.
func Gen() error {
	fset := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fset, "dummy", nil, parser.Mode(0))
	funcs := []domain.Func{}
	testFuncs := domain.Funcs{Names: map[string]map[string]domain.Func{}}
	testFileNameSet := set.NewStrSet()
	for pkgName, pkg := range pkgs {
		for fileName, f := range pkg.Files {
			if strings.HasSuffix(fileName, "_test.go") {
				testFileNameSet.Add(fileName)
			}
			ast.Inspect(f, func(n ast.Node) bool {
				if ident, ok := n.(*ast.FuncDecl); ok {
					fName := ident.Name.Name
					if !isExposed(fName) {
						return true
					}
					structName := ""
					if ident.Recv != nil {
						structName = fmt.Sprintf("%v", ident.Recv.List[0].Type)
					}
					fnc := domain.Func{
						Name:        fName,
						FileName:    fileName,
						StructName:  structName,
						PackageName: pkgName,
					}
					if fnc.IsTestFile() {
						if fnc.IsTest() {
							testFuncs.Add(fnc)
						}
						return true
					}
					funcs = append(funcs, fnc)
				}
				return true
			})
		}
	}
	// file name -> appended test code
	addedCodes := map[string]string{}
	for _, f := range funcs {
		if _, ok := testFuncs.Names[f.TestFuncName()]; ok {
			continue
		}
		if m, ok := addedCodes[f.TestFileName()]; ok {
			m += f.TestCode()
			addedCodes[f.TestFileName()] = m
			continue
		}
		if testFileNameSet.Has(f.TestFileName()) {
			// append
			addedCodes[f.TestFileName()] = f.TestCode()
			continue
		}
		// create
		addedCodes[f.TestFileName()] = fmt.Sprintf(`package %s

import (
	"testing"
)
%s`, f.PackageName, f.TestCode())
	}

	for fName, code := range addedCodes {
		if testFileNameSet.Has(fName) {
			// append
			f, err := os.OpenFile(fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to open file: %s", fName))
			}
			defer f.Close()
			fmt.Printf("add a test skelton code to a test file: %s\n", fName)
			if _, err := f.Write([]byte(code)); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to write file: %s", fName))
			}
			continue
		}
		// create
		fmt.Printf("create a test file: %s\n", fName)
		if err := ioutil.WriteFile(fName, []byte(code), 0644); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to create test file: %s", fName))
		}
	}
	return nil
}
