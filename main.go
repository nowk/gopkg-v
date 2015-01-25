package main

import (
	"fmt"
	"gopkg.in/nowk/gopkg-v.v0/pkg"
	"os"
)

const (
	VERSION = "0.0.0"
)

var usage = `---
Usage:
  gopkg-v <user/package> --gopkg-version <number>
  gopkg-v -h | --help
  gopkg-v -v | --version

Options:
  --package <user/package>  package name to version (user/package format)
  --gopkg-version <number>  gopkg.in version
                            If versions exist and none is provided, it will 
                            automatically increment to the next version.
                            If there are no versions, it will create a .v0

  -h, --help                output help information
  -v, --version             output version
`

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	cfg, err := pkg.ParseArgs(os.Args)
	check(err)

	p, err := pkg.Open(cfg)
	check(err)

	v, err := p.NewVersion()
	check(err)
	fmt.Fprintf(os.Stdout, "%s/%s is now at version.%d\n", p.User, p.Name, v.Version)

	// log.Print("go get ", goget)
	// err = exec.Command("go", "get", goget).Run()
	// if err != nil {
	// 	log.Fatalf("fatal: could not go get %s\n%s", goget, err)
	// }
}
