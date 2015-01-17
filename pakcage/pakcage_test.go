package pakcage

import (
	"gopkg.in/nowk/assert.v2"
	"os"
	"path"
	"testing"
)

var (
	pwd      string
	srcGopkg string
)

func init() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	pwd = path.Join(pwd, "test")
}

func setup() func() {
	origGoPath := GoPath
	GoPath = pwd

	srcGopkg = "src/gopkg.in/foo"

	os.RemoveAll(path.Join(GoPath, srcGopkg))
	os.MkdirAll(path.Join(GoPath, srcGopkg, "awesome.v2"), os.ModePerm)
	os.Create(path.Join(GoPath, srcGopkg, "awesome.v3"))
	// os.Create(path.Join(GoPath, srcGopkg, "more_awesome.v3"))

	return func() {
		GoPath = origGoPath
	}
}

func TestRepoPaths(t *testing.T) {
	teardown := setup()
	defer teardown()

	p, err := New("foo", "awesome")
	if err != nil {
		t.Fatal(err)
	}

	gopk := path.Join(pwd, srcGopkg)
	gith := path.Join(pwd, "src", "github.com", "foo")

	assert.Equal(t, gopk, p.Gopkg)
	assert.Equal(t, gith, p.Github)
}

func TestGetLatestVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	p, err := New("foo", "awesome")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 3, p.Version)
}

func TestBumpsVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	p, err := New("foo", "awesome")
	if err != nil {
		t.Fatal(err)
	}

	err = p.NewVersion()
	if err != nil {
		t.Fatal(err)
	}

	str, err := os.Readlink(path.Join(p.Gopkg, "awesome.v4"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, path.Join(p.Github, "awesome"), str)
	assert.Equal(t, 4, p.Version)
}

func TestCreatesNewVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	p, err := New("foo", "more_awesome")
	if err != nil {
		t.Fatal(err)
	}

	err = p.NewVersion()
	if err != nil {
		t.Fatal(err)
	}

	str, err := os.Readlink(path.Join(p.Gopkg, "more_awesome.v0"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, path.Join(p.Github, "more_awesome"), str)
	assert.Equal(t, 0, p.Version)
}

func TestCreateSpecificVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	p, err := New("foo", "more_awesome")
	if err != nil {
		t.Fatal(err)
	}

	err = p.NewVersion(9)
	if err != nil {
		t.Fatal(err)
	}

	str, err := os.Readlink(path.Join(p.Gopkg, "more_awesome.v9"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, path.Join(p.Github, "more_awesome"), str)
	assert.Equal(t, 9, p.Version)
}

// TODO go gets previous version
// TODO does not get v0 when going to v1
