package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
)

const (
	GITHUB_PATH = "src/github.com"
	GOPKG_PATH  = "src/gopkg.in"
)

var (
	GoPath = os.Getenv("GOPATH")
)

type version struct {
	// Path is the path to $GOPATH/src/gopkg.in/user/package.vN
	Path string

	// Name is the versioned named package.vN
	Name    string
	Version int

	// WorkingDir is true if the version is liked to the github.com/user/package
	WorkingDir bool
}

type Pkg struct {
	User string
	Name string

	// Source is the path to $GOPATH/src/github.com/user/package
	Source string

	// Dest is the root gopkg.in path $GOPATH/src/gopkg.in/user
	Dest string

	// Versions is always assumed to be in v0, v1, v2, v3..., vN order
	Versions []*version
}

func Open(conf *Config) (*Pkg, error) {
	pkg := &Pkg{
		User:   conf.User,
		Name:   conf.Repo,
		Source: path.Join(GoPath, GITHUB_PATH, conf.User, conf.Repo),
		Dest:   path.Join(GoPath, GOPKG_PATH, conf.User),
	}

	_, err := os.Open(pkg.Source)
	if err != nil {
		return nil, fmt.Errorf("package %s/%s: no such package", pkg.User,
			pkg.Name)
	}

	if err := queryVersions(pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}

func isWorkingDir(strPath string) bool {
	// TODO distinguish between non-symlink and bad paths
	_, err := os.Readlink(strPath)
	return err == nil
}

// CurrentVersion returns the last version in the versions list
func (p Pkg) CurrentVersion() *version {
	l := len(p.Versions)
	if l > 0 {
		return p.Versions[l-1]
	}

	return nil
}

// unlink calls the unlink command to ensure an "unlink" occurs vs doing a
// straight rm or rmdir
func unlink(v *version) error {
	if v == nil {
		return nil
	}

	// check if file exists
	if _, err := os.Readlink(v.Path); err != nil {
		return nil
	}

	return exec.Command("unlink", v.Path).Run()
}

func nextVersion(v *version) int {
	if v != nil {
		return v.Version + 1
	}

	return 0
}

// NewVersionAt creates a new version at a particular version N.
// WARNING this does not remove any existing symlinks to the source directory
func (p *Pkg) NewVersionAt(n int) (*version, error) {
	cur := p.CurrentVersion()
	if cur != nil && cur.Version == n {
		return nil, fmt.Errorf("package %s/%s: is already at version %d", p.User,
			p.Name, cur.Version)
	}

	dir := fmt.Sprintf("%s.v%d", p.Name, n)
	sym := path.Join(p.Dest, dir)
	if err := os.Symlink(p.Source, sym); err != nil {
		return nil, err
	}

	return p.addVersion(dir)
}

// NewVersion creates a new version at the next version N. It will also unlink
// the current version so only one active symlink to the source directory
func (p *Pkg) NewVersion() (*version, error) {
	cur := p.CurrentVersion()

	v, err := p.NewVersionAt(nextVersion(cur))
	if err != nil {
		return v, err
	}

	if err := unlink(cur); err != nil {
		return v, err
	}

	return v, nil
}

func (p *Pkg) addVersion(name string) (*version, error) {
	reg, err := regexp.Compile(`\d+$`)
	if err != nil {
		return nil, err
	}

	nStr := reg.FindString(name)
	n, err := strconv.ParseInt(nStr, 10, 64)
	if err != nil {
		return nil, err
	}

	vPath := path.Join(p.Dest, name)
	v := &version{
		Path:       vPath,
		Name:       name,
		Version:    int(n),
		WorkingDir: isWorkingDir(vPath),
	}
	p.Versions = append(p.Versions, v)
	return v, nil
}

func queryVersions(pkg *Pkg) error {
	reg, err := regexp.Compile(fmt.Sprintf(`^%s.v\d+$`, pkg.Name))
	if err != nil {
		return err
	}

	f, err := os.Open(pkg.Dest)
	if err != nil {
		return err
	}
	dirs, err := f.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, v := range dirs {
		if !reg.MatchString(v) {
			continue
		}

		if _, err := pkg.addVersion(v); err != nil {
			return err
		}
	}

	return nil
}