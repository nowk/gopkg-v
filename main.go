package main

import (
	"flag"
	"fmt"
	"gopkg.in/nowk/gopkg-v.v0/pakcage"
	// "github.com/nowk/gopkg-v/pakcage"
	"log"
	"os"
	"strings"
)

const (
	VERSION = "0.0.0"
)

var usage = `

---
Usage:
  gopkg-v --package <user/package> [--gopkg-version <number>]
  gopkg-v -h | --help
  gopkg-v -v | --version

Options:
  --package <user/package>  package name to version (user/package format)
  --gopkg-version <number>  gopkg.in version [default: -1] 
                            If versions exist and none is provided, it will 
                            automatically increment to the next version.
                            If there are no versions, it will create a .v0

  -h, --help                output help information
  -v, --version             output version
`

var (
	packagename    string
	packageversion int

	version bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gopkg-v -- Version %s%s\n", VERSION, usage)
		os.Exit(0)
	}

	flag.StringVar(&packagename, "package", "", "")
	flag.IntVar(&packageversion, "gopkg-version", -1, "")

	flag.BoolVar(&version, "version", false, "")
	flag.BoolVar(&version, "v", version, "")

	flag.Parse()

	if version {
		fmt.Fprintf(os.Stderr, "gopkg-v -- Version %s\n", VERSION)
		os.Exit(0)
	}

	if packagename == "" {
		flag.Usage()
	}
}

func main() {
	s := strings.Split(packagename, "/")
	if len(s) < 2 {
		log.Fatal("error: package must be in a <user>/<package> format")
	}
	user, pack := s[0], s[1]

	p, err := pakcage.New(user, pack)
	if err != nil {
		log.Fatal(err)
	}
	if err := p.NewVersion(packageversion); err != nil {
		log.Fatal(err)
	}

	log.Printf("%s is now at version.%d", packagename, p.Version)

	// log.Print("go get ", goget)
	// err = exec.Command("go", "get", goget).Run()
	// if err != nil {
	// 	log.Fatalf("fatal: could not go get %s\n%s", goget, err)
	// }
}
