// Package main is an entrypoint
package main

import (
	"flag"
	"fmt"
)

var versionFlag = flag.Bool("version", false, "print version")

var version string

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("migr8 version: %s\n", version)
		return
	}
}
