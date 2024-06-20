package linter

import (
	"os"
)

type CheckPos struct {
	Coll int32
	Line int32
}

type CheckError struct {
	Code  string
	Level string
	Msg   string
}

type CheckResult struct {
	Pos CheckPos
	Err CheckError
}

type ILinter interface {
	Check(file *os.File) ([]CheckResult, error)
}

type HadolintResult struct {
	Code    string `json:"code"`
	Coll    int32  `json:"column"`
	File    string `json:"file,omitempty"`
	Level   string `json:"level"`
	Line    int32  `json:"line"`
	Message string `json:"message"`
}
