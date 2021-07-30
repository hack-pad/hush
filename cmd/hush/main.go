package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hack-pad/hush"
)

func main() {
	log.SetOutput(ioutil.Discard)
	os.Exit(hush.Run())
}
