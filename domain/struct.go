package domain

import (
	"fmt"
	"path/filepath"
	"strings"
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
	return fmt.Sprintf("Test%s", f.Name)
}

// TestCode returns the boilerplate of test function.
func (f Func) TestCode() string {
	return fmt.Sprintf(`
func Test%s%s(t *testing.T) {
}
`, f.StructName, f.Name)
}

// IsTest returns whether the function is test function or not.
func (f Func) IsTest() bool {
	return strings.HasPrefix(f.Name, "Test")
}

// IsTestFile returns whether the function is written in a test file or not.
func (f Func) IsTestFile() bool {
	return strings.HasSuffix(f.FileName, "_test.go")
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
