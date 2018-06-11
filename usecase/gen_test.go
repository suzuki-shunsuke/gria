package usecase

import (
	"fmt"
	"testing"

	"github.com/suzuki-shunsuke/go-set"

	"github.com/suzuki-shunsuke/gria/domain"
)

func TestGetFuncs(t *testing.T) {
}

func TestGetCodes(t *testing.T) {
	funcs := []domain.Func{{
		Name:        "Append",
		StructName:  "FileWriter",
		FileName:    "./infra/file.go",
		PackageName: "infra",
	}}
	testFuncs := domain.CreateFuncs()
	testFuncs.Add(domain.Func{
		Name:        "TestFileWriterAppend",
		StructName:  "",
		FileName:    "./infra/file_test.go",
		PackageName: "infra",
	})
	testFileNameSet := set.NewStrSet("./infra/file_test.go")
	ac := GetCodes(funcs, testFuncs, testFileNameSet)
	if len(ac.Codes) != 0 {
		fmt.Println(ac.Codes)
		t.Fatalf(`len(ac.Codes) = %d, wanted 0`, len(ac.Codes))
	}
}

func TestWriteCodes(t *testing.T) {
}

func TestGen(t *testing.T) {
}
