package shuai_utils

import (
	"fmt"
	"testing"
)

func TestPanic(t *testing.T) {
	var err error
	Panic(err)
	err = fmt.Errorf("err")
	Panic(err)
}
