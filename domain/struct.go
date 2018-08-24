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
	// function name -> package name -> Func
	Funcs map[string]map[string]Func

	// AddedCodes represents added codes.
	AddedCodes struct {
		// file name -> appended test code
		Codes           map[string]string
		TestFileNameSet set.StrSet
	}

	// GenArgs is arguments of usecase.Gen function.
	GenArgs struct {
		Targets    []string
		FileWriter FileWriter
	}
)

// TestFileName returns a test file name.
func (f *Func) TestFileName() string {
	return fmt.Sprintf(
		"%s_test.go",
		f.FileName[:len(f.FileName)-len(filepath.Ext(f.FileName))])
}

// TestFuncName returns a test function name.
func (f *Func) TestFuncName() string {
	a := fmt.Sprintf("%s%s", f.StructName, f.Name)
	if a != strings.Title(a) {
		a = fmt.Sprintf("_%s", a)
	}
	return fmt.Sprintf("Test%s", a)
}

// TestCode returns the boilerplate of test function.
func (f *Func) TestCode() string {
	return fmt.Sprintf(`
func %s(t *testing.T) {
}
`, f.TestFuncName())
}

// IsTest returns whether the function is test function or not.
func (f *Func) IsTest() bool {
	return f.IsTestFile() && strings.HasPrefix(f.Name, "Test")
}

// IsTestFile returns whether the function is written in a test file or not.
func (f *Func) IsTestFile() bool {
	return strings.HasSuffix(f.FileName, "_test.go")
}

// Add adds a function.
func (funcs Funcs) Add(f Func) error {
	if funcs == nil {
		return fmt.Errorf("Funcs is nil")
	}
	if m, ok := funcs[f.Name]; ok {
		m[f.PackageName] = f
		funcs[f.Name] = m
		return nil
	}
	funcs[f.Name] = map[string]Func{f.PackageName: f}
	return nil
}

// HasTest returns whether the test has already existed.
func (funcs Funcs) HasTest(f Func) bool {
	if funcs == nil {
		return false
	}
	pkgs, ok := funcs[f.TestFuncName()]
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
func CreateAddedCodes(testFileNameSet set.StrSet) AddedCodes {
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
