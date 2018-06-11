package domain

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/go-set"
)

type (
	// Func represents a function code.
	Func struct {
		Name        string
		StructName  string
		FileName    string
		PackageName string
	}

	// Funcs represents functions.
	Funcs struct {
		// func name -> package name -> func
		Names map[string]map[string]Func
	}

	// AddedCodes represents added codes.
	AddedCodes struct {
		// file name -> appended test code
		Codes           map[string]string
		TestFileNameSet *set.StrSet
	}

	// GenArgs is arguments of usecase.Gen function.
	GenArgs struct {
		Targets    []string
		FileWriter FileWriter
	}
)

// TestFileName returns a test file name.
func (f Func) TestFileName() string {
	return fmt.Sprintf(
		"%s_test.go",
		f.FileName[:len(f.FileName)-len(filepath.Ext(f.FileName))])
}

// TestFuncName returns a test function name.
func (f Func) TestFuncName() string {
	return fmt.Sprintf("Test%s%s", f.StructName, f.Name)
}

// TestCode returns the boilerplate of test function.
func (f Func) TestCode() string {
	return fmt.Sprintf(`
func %s(t *testing.T) {
}
`, f.TestFuncName())
}

// IsTest returns whether the function is test function or not.
func (f Func) IsTest() bool {
	return strings.HasPrefix(f.Name, "Test")
}

// IsTestFile returns whether the function is written in a test file or not.
func (f Func) IsTestFile() bool {
	return strings.HasSuffix(f.FileName, "_test.go")
}

// CreateFuncs returns a Funcs.
func CreateFuncs() Funcs {
	return Funcs{Names: map[string]map[string]Func{}}
}

// Add adds a function.
func (funcs *Funcs) Add(f Func) {
	if m, ok := funcs.Names[f.Name]; ok {
		m[f.PackageName] = f
		funcs.Names[f.Name] = m
		return
	}
	funcs.Names[f.Name] = map[string]Func{f.PackageName: f}
}

// HasTest returns whether the test has already existed.
func (funcs Funcs) HasTest(f Func) bool {
	if funcs.Names == nil {
		return false
	}
	pkgs, ok := funcs.Names[f.TestFuncName()]
	if !ok {
		return false
	}
	if pkgs == nil {
		return false
	}
	_, ok = pkgs[f.PackageName]
	return ok
}

// CreateAddedCodes returns an AddedCodes.
func CreateAddedCodes(testFileNameSet *set.StrSet) AddedCodes {
	return AddedCodes{
		Codes:           map[string]string{},
		TestFileNameSet: testFileNameSet}
}

// Init initialize an AddedCodes.
func (ac *AddedCodes) Init() {
	if ac.Codes == nil {
		ac.Codes = map[string]string{}
	}
}

// Add adds a test code to an AddedCodes.
func (ac *AddedCodes) Add(f Func) {
	ac.Init()
	if m, ok := ac.Codes[f.TestFileName()]; ok {
		m += f.TestCode()
		ac.Codes[f.TestFileName()] = m
		return
	}
	if ac.TestFileNameSet.Has(f.TestFileName()) {
		// append
		ac.Codes[f.TestFileName()] = f.TestCode()
		return
	}
	// create
	ac.Codes[f.TestFileName()] = fmt.Sprintf(`package %s

import (
	"testing"
)
%s`, f.PackageName, f.TestCode())
}
