package gogit

import (
	"bytes"
	"fmt"
)

type ExecErr struct {
	Msg    string
	Stdout bytes.Buffer
	Stderr bytes.Buffer
	Base   error
}

func execErr(msg string, stdout bytes.Buffer, stderr bytes.Buffer, base error) *ExecErr {
	return &ExecErr{
		Msg:    msg,
		Stdout: stdout,
		Stderr: stderr,
		Base:   base,
	}
}

func (e *ExecErr) Error() string {
	return fmt.Sprintf("Exec error: %s: %s", e.Msg, e.Base)
}

func (e *ExecErr) Unwrap() error {
	return e.Base
}

var _ error = &ExecErr{}
