package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

type (
	Func struct {
		Name        string
		StructName  string
		FileName    string
		PackageName string
	}
	Funcs struct {
		// func name -> package name -> func
		Names map[string]map[string]Func
	}
	// GenArgs is arguments of usecase.Gen function.
	GenArgs struct {
		Targets []string
	}
)

func (f Func) TestFileName() string {
	return fmt.Sprintf(
		"%s_test.go",
		f.FileName[:len(f.FileName)-len(filepath.Ext(f.FileName))])
}

func (f Func) TestFuncName() string {
	return fmt.Sprintf("Test%s", f.Name)
}

func (f Func) TestCode() string {
	return fmt.Sprintf(`
func Test%s%s(t *testing.T) {
}
`, f.StructName, f.Name)
}

func (f Func) IsTest() bool {
	return strings.HasPrefix(f.Name, "Test")
}

func (f Func) IsTestFile() bool {
	return strings.HasSuffix(f.FileName, "_test.go")
}

func (funcs *Funcs) Add(f Func) {
	if m, ok := funcs.Names[f.Name]; ok {
		m[f.PackageName] = f
		funcs.Names[f.Name] = m
		return
	}
	funcs.Names[f.Name] = map[string]Func{f.PackageName: f}
}
