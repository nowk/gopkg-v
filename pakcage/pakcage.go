package pakcage

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

const (
	GOPKG_HOST  = "gopkg.in"
	GITHUB_HOST = "github.com"
)

var (
	GoPath = os.Getenv("GOPATH")
)

func init() {
	if GoPath == "" {
		panic("GOPATH is not set")
	}
}

type Pakcage struct {
	user    string
	name    string
	Gopkg   string
	Github  string
	Version int // eg. 1, 2, etc...
}

func parseVersion(str string, reg *regexp.Regexp) (int, error) {
	n := reg.FindAllStringSubmatch(str, -1)[0][1]
	v, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		return -1, err
	}

	return int(v), nil
}

func currentVersion(p *Pakcage) (int, error) {
	f, err := os.Open(p.Gopkg)
	if err != nil {
		return -1, err
	}

	reg := regexp.MustCompile(fmt.Sprintf(`^%s\.v(\d+)`, p.name))
	d, err := matchDirs(f, reg)
	if err != nil {
		return -1, err
	}
	if last, ok := d.Last(); ok {
		return parseVersion(last, reg)
	}

	return -1, nil
}

// New returns a preconfigured Pakcage with the latest version in the gopkg set.
// Unversionable Pakcages will be set at version -1
func New(user, name string) (*Pakcage, error) {
	p := &Pakcage{
		user:   user,
		name:   name,
		Gopkg:  path.Join(GoPath, "src", GOPKG_HOST, user),
		Github: path.Join(GoPath, "src", GITHUB_HOST, user),
	}

	v, err := currentVersion(p)
	if err != nil {
		return nil, err
	}

	p.Version = v

	return p, nil
}

func (p Pakcage) gopkgVerURL(v int) string {
	return path.Join(p.Gopkg, fmt.Sprintf("%s.v%d", p.name, v))
}

// NewVersion unlinks the current version, links the new version and bumps the
// Version value of this Pakcage
func (p *Pakcage) NewVersion(i ...int) error {
	v := -1
	if len(i) > 0 {
		v = i[0]
	}

	// no input version or < 0, use the parsed Pakcage version
	if v < 0 {
		v = p.Version + 1
	}

	err := p.unlink()
	if err != nil {
		return err
	}

	err = p.linkVersion(v)
	if err != nil {
		return err
	}

	p.Version = v // update version

	return nil
}

// linkVersion links the current gopkg version to the github repo
func (p Pakcage) linkVersion(v int) error {
	return os.Symlink(path.Join(p.Github, p.name), p.gopkgVerURL(v))
}

// unlink unlinks the current version. It will error if path to current version
// is not a link (is a directory). It will not attempt to unlink a .v-1 version.
func (p Pakcage) unlink() error {
	if p.Version < 0 {
		return nil
	}

	return exec.Command("unlink", p.gopkgVerURL(p.Version)).Run()
}
