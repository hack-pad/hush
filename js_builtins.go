// +build js

package hush

import (
	"fmt"
	"io/ioutil"
	"strings"
	"syscall/js"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

var (
	jsFunction = js.Global().Get("Function")
)

func init() {
	builtins["jseval"] = jseval
	builtins["jsdownload"] = jsdownload
	color.NoColor = false // override, since wasm isn't considered a "tty"
}

func jsEval(funcStr string, args ...interface{}) js.Value {
	f := jsFunction.Invoke(`"use strict";` + funcStr)
	return f.Invoke(args...)
}

func jseval(term console, args ...string) error {
	if len(args) < 1 {
		return errors.New("Must provide a string to run as a function")
	}
	result := jsEval(args[0], strings.Join(args[1:], " "))
	fmt.Fprintln(term.Stdout(), result)
	return nil
}

func jsdownload(term console, args ...string) error {
	if len(args) < 1 {
		return errors.New("Must provide a file to download")
	}
	filePath := args[0]
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "Error reading file for download")
	}
	startDownload("", filePath, fileContents)
	return nil
}
