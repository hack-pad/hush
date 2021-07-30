package hush

import "fmt"

type exitErr struct {
	Code int
}

func (e *exitErr) Error() string {
	return fmt.Sprintf("exit code %d", e.Code)
}
