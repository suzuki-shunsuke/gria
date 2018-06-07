package usecase

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/suzuki-shunsuke/go-set"

	"github.com/suzuki-shunsuke/gria/domain"
)

// GetFuncs extracts target functions and existing test functions and test file paths from given packages.
func GetFuncs(pkgs map[string]*ast.Package) ([]domain.Func, domain.Funcs, *set.StrSet) {
	funcs := []domain.Func{}
	testFuncs := domain.Funcs{Names: map[string]map[string]domain.Func{}}
	testFileNameSet := set.NewStrSet()
	for pkgName, pkg := range pkgs {
		for fileName, f := range pkg.Files {
			// judge test file by file name
			if strings.HasSuffix(fileName, "_test.go") {
				testFileNameSet.Add(fileName)
			}
			ast.Inspect(f, func(n ast.Node) bool {
				// filter function
				if ident, ok := n.(*ast.FuncDecl); ok {
					fName := ident.Name.Name
					// get struct name
					structName := ""
					if ident.Recv != nil {
						expr := ident.Recv.List[0].Type
						if p, ok := expr.(*ast.StarExpr); ok {
							structName = fmt.Sprintf("%v", p.X)
						} else {
							structName = fmt.Sprintf("%v", expr)
						}
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
	return funcs, testFuncs, testFileNameSet
}

// GetCodes returns added test codes to each test files.
func GetCodes(funcs []domain.Func, testFuncs domain.Funcs, testFileNameSet *set.StrSet) map[string]string {
	// file name -> appended test code
	addedCodes := map[string]string{}
	// generate test codes
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
	return addedCodes
}

// WriteCodes writes test codes to test files.
func WriteCodes(addedCodes map[string]string, testFileNameSet *set.StrSet, fileWriter domain.FileWriter) error {
	for fName, code := range addedCodes {
		if testFileNameSet.Has(fName) {
			// append
			if err := fileWriter.Append(fName, []byte(code)); err != nil {
				return err
			}
			continue
		}
		// create
		fmt.Printf("create a test file: %s\n", fName)
		if err := fileWriter.Create(fName, []byte(code)); err != nil {
			return err
		}
	}
	return nil
}

// Gen generates test scaffold.
func Gen(args domain.GenArgs) error {
	if len(args.Targets) == 0 {
		return nil
	}
	for _, arg := range args.Targets {
		if filepath.Ext(arg) == ".go" {
			return fmt.Errorf("argument must not be file but directory: %s", arg)
		}
	}
	fset := token.NewFileSet()
	for _, p := range args.Targets {
		pkgs, err := parser.ParseDir(fset, p, nil, parser.Mode(0))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse package: %s", p))
		}
		funcs, testFuncs, testFileNameSet := GetFuncs(pkgs)

		addedCodes := GetCodes(funcs, testFuncs, testFileNameSet)
		if err := WriteCodes(addedCodes, testFileNameSet, args.FileWriter); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to write test files of package %s", p))
		}
	}
	return nil
}
