package hush

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Run runs the hush shell
func Run() int {
	cancel, err := ttySetup()
	if err != nil {
		panic(err)
	}
	defer cancel()
	return run(os.Stdin, os.Stdout, os.Stderr, os.Args)
}

func run(in io.Reader, out, outErr io.Writer, args []string) int {
	set := flag.NewFlagSet(args[0], flag.ContinueOnError)
	command := set.String("c", "", "Read and execute commands from the given string value.")
	err := set.Parse(args[1:])
	if err != nil {
		fmt.Fprintln(outErr, err)
		return 2
	}

	var reader io.RuneReader
	if *command != "" {
		reader = newRuneReader(strings.NewReader(*command))
	} else {
		reader = newRuneReader(in)
	}
	out, err = newCarriageReturnWriter(out)
	if err != nil {
		panic(err)
	}
	outErr, err = newCarriageReturnWriter(outErr)
	if err != nil {
		panic(err)
	}
	term := newTerminal(out, outErr)

	return term.ReadEvalPrintLoop(reader)
}
