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
func GetFuncs(pkgPath string) ([]domain.Func, domain.Funcs, set.StrSet, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, parser.Mode(0))
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("failed to parse package: %s", pkgPath))
	}
	funcs := []domain.Func{}
	testFuncs := domain.Funcs{}
	testFileNameSet := set.NewStrSet()
	for pkgName, pkg := range pkgs {
		// ignore main package
		if pkgName == "main" {
			continue
		}
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
					// ignore init functions
					if fnc.StructName == "" && fnc.Name == "init" {
						return true
					}
					funcs = append(funcs, fnc)
				}
				return true
			})
		}
	}
	return funcs, testFuncs, testFileNameSet, nil
}

// GetCodes returns added test codes to each test files.
func GetCodes(funcs []domain.Func, testFuncs domain.Funcs, testFileNameSet set.StrSet) domain.AddedCodes {
	// file name -> appended test code
	addedCodes := domain.CreateAddedCodes(testFileNameSet)
	// generate test codes
	for _, f := range funcs {
		if testFuncs.HasTest(f) {
			continue
		}
		addedCodes.Add(f)
	}
	return addedCodes
}

// WriteCodes writes test codes to test files.
func WriteCodes(addedCodes domain.AddedCodes, testFileNameSet set.StrSet, fileWriter domain.FileWriter) error {
	for fName, code := range addedCodes.Codes {
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
	targets := set.NewStrSet(args.Targets...)
	for arg := range targets {
		if filepath.Ext(arg) == ".go" {
			return fmt.Errorf("argument must not be file but directory: %s", arg)
		}
	}
	for p := range targets {
		funcs, testFuncs, testFileNameSet, err := GetFuncs(p)
		if err != nil {
			return err
		}
		addedCodes := GetCodes(funcs, testFuncs, testFileNameSet)
		if err := WriteCodes(addedCodes, testFileNameSet, args.FileWriter); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to write test files of package %s", p))
		}
	}
	return nil
}
