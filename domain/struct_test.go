package domain

import (
	"testing"
)

func TestFuncTestFileName(t *testing.T) {
	f := Func{FileName: "foo.go"}
	act := f.TestFileName()
	exp := "foo_test.go"
	if act != exp {
		t.Fatalf("%s != %s", act, exp)
	}
}

func TestFuncTestFuncName(t *testing.T) {
	f := Func{Name: "Foo"}
	act := f.TestFuncName()
	exp := "TestFoo"
	if act != exp {
		t.Fatalf("%s != %s", act, exp)
	}

	f = Func{Name: "Foo", StructName: "Bar"}
	act = f.TestFuncName()
	exp = "TestBarFoo"
	if act != exp {
		t.Fatalf("%s != %s", act, exp)
	}
}

func TestFuncTestCode(t *testing.T) {
	f := Func{Name: "Foo"}
	act := f.TestCode()
	exp := `
func TestFoo(t *testing.T) {
}
`
	if act != exp {
		t.Fatalf("%s != %s", act, exp)
	}
}

func TestFuncIsTest(t *testing.T) {
	f := Func{Name: "Foo"}
	if f.IsTest() {
		t.Fatalf("Foo() is not test function")
	}
}

func TestFuncIsTestFile(t *testing.T) {
	f := Func{FileName: "foo.go"}
	if f.IsTestFile() {
		t.Fatalf("foo.go is not test function")
	}
}

func TestFuncsAdd(t *testing.T) {
	f := Func{PackageName: "p", Name: "Foo", FileName: "foo.go"}
	funcs := CreateFuncs()
	funcs.Add(f)
	funcs = Funcs{Names: map[string]map[string]Func{"Foo": {}}}
	funcs.Add(f)
}

func TestCreateFuncs(t *testing.T) {
	CreateFuncs()
}

func TestFuncsHasTest(t *testing.T) {
	funcs := CreateFuncs()
	f := Func{PackageName: "p", Name: "Foo", FileName: "foo.go"}
	if funcs.HasTest(f) {
		t.Fatal("funcs.HasTest(f) = true, wanted false")
	}
	ft := Func{PackageName: "p", Name: "TestFoo", FileName: "foo_test.go"}
	funcs.Add(ft)
	if !funcs.HasTest(f) {
		t.Fatal("funcs.HasTest(f) = false, wanted true")
	}
}

func TestCreateAddedCodes(t *testing.T) {
	CreateAddedCodes(nil)
}

func TestAddedCodesInit(t *testing.T) {
	ac := CreateAddedCodes(nil)
	ac.Init()
}

func TestAddedCodesAdd(t *testing.T) {
	ac := CreateAddedCodes(nil)
	f := Func{PackageName: "p", Name: "Foo", FileName: "foo.go"}
	ac.Add(f)
}
