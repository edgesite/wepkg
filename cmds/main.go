package main

import (
	"github.com/edgesite/wepkg"
	"os"
)

func main() {
	os.Exit(wepkg.Unpack(os.Args[1], true))
}
