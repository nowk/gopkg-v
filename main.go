package main

import (
	"flag"
	"fmt"
	"gopkg.in/nowk/gopkg-v.v0/pkg"
	"os"
)

const (
	VERSION = "0.0.0"
)

var usage = `---
Usage:
  gopkg-v <user/package> --new-version [<number>]
  gopkg-v --help
  gopkg-v --version

Options:
  <user/package>            package name to version (eg. nowk/gopkg-v)
  --new-version [<number>]  creates a version at the given number. If number is
                            blank will create at the current version + 1. A
                            value of 0 will behave like a blank value if
                            previous versions exist.
                            ex: --new-version 1 -> package.v1

  --help                    output usage information
  --version                 output version
`

func printusage() {
	fmt.Fprintf(os.Stdout, "%s\n", usage)
}

func init() {
	hlp := flag.Bool("help", false, "output usage information")
	ver := flag.Bool("version", false, "output version")
	flag.Parse()

	if *hlp {
		printusage()
		os.Exit(0)
	}

	if *ver {
		fmt.Fprintf(os.Stdout, "gopkg-v version %s\n", VERSION)
		os.Exit(0)
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		_, ok := err.(*pkg.ArgError)
		if ok {
			printusage()
		}
		os.Exit(1)
	}
}

func main() {
	cfg, err := pkg.ParseArgs(os.Args)
	check(err)

	p, err := pkg.Open(cfg)
	check(err)

	v, err := p.NewVersion(cfg.Version)
	check(err)
	fmt.Fprintf(os.Stdout, "%s/%s is now at version.%d\n", p.User, p.Name, v.Version)

	// log.Print("go get ", goget)
	// err = exec.Command("go", "get", goget).Run()
	// if err != nil {
	// 	log.Fatalf("fatal: could not go get %s\n%s", goget, err)
	// }
}
