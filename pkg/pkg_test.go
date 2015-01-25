package pkg

import (
	"fmt"
	"gopkg.in/nowk/gopkg-v.v0/pkg/testing/assert"
	"os"
	"path"
	"testing"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	GoPath = fmt.Sprintf("%s/test_gopath", pwd)
}

func setup() func() {
	var (
		gopkg  = path.Join(GoPath, "src/gopkg.in")
		github = path.Join(GoPath, "src/github.com")
	)

	githubFoo := path.Join(github, "foo")
	gopkgFoo := path.Join(gopkg, "foo")

	os.MkdirAll(path.Join(gopkgFoo, "awesome.v2"), os.ModePerm)
	os.MkdirAll(path.Join(gopkgFoo, "awesome.v3"), os.ModePerm)
	os.Symlink(path.Join(githubFoo, "awesome"), path.Join(gopkgFoo, "awesome.v4"))

	return func() {
		os.RemoveAll(gopkgFoo)
	}
}

func TestReturnsAParsedPackage(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "awesome",
	})
	assert.Nil(t, err)
	assert.Equal(t, "foo", pkg.User)
	assert.Equal(t, "awesome", pkg.Name)
	assert.Equal(t, fmt.Sprintf("%s/src/github.com/foo/awesome", GoPath),
		pkg.Source)
	assert.Equal(t, fmt.Sprintf("%s/src/gopkg.in/foo", GoPath), pkg.Dest)
	assert.Equal(t, 3, len(pkg.Versions))

	for i, v := range []struct {
		Name       string
		Version    int
		WorkingDir bool
	}{
		{"awesome.v2", 2, false},
		{"awesome.v3", 3, false},
		{"awesome.v4", 4, true},
	} {
		assert.Equal(t, v.Name, pkg.Versions[i].Name)
		assert.Equal(t, v.Version, pkg.Versions[i].Version)
		assert.Equal(t, v.WorkingDir, pkg.Versions[i].WorkingDir)
	}
}

func TestNewVersionCreatesNewLinkToWorkingDir(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "awesome",
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := pkg.NewVersion(0)
	assert.Nil(t, err)
	assert.Equal(t, "awesome.v5", v.Name)
	assert.Equal(t, 5, v.Version)
	assert.True(t, v.WorkingDir)
	assert.Symlink(t, pkg.Source, v.Path)

	cur := pkg.CurrentVersion()
	assert.Equal(t, v, cur)
}

func TestNewVersionUnlinksTheOldVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "awesome",
	})
	if err != nil {
		t.Fatal(err)
	}
	cur := pkg.CurrentVersion()
	if _, err := pkg.NewVersion(0); err != nil {
		t.Fatal(err)
	}
	assert.Unlink(t, cur.Path)
}

func TestNoVersionStartsAt0(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "more_awesome",
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := pkg.NewVersion(0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "more_awesome.v0", v.Name)
	assert.Equal(t, 0, v.Version)
	assert.True(t, v.WorkingDir)
	assert.Symlink(t, pkg.Source, v.Path)
}

func TestNoSourceRepoExist(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "bar",
	})
	assert.Nil(t, pkg)
	assert.Equal(t, "package foo/bar: no such package", err.Error())
}

func TestCreateVersionSpecificVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "awesome",
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := pkg.NewVersion(5)
	assert.Nil(t, err)
	assert.Equal(t, "awesome.v5", v.Name)
	assert.Equal(t, 5, v.Version)
	assert.True(t, v.WorkingDir)
	assert.Symlink(t, pkg.Source, v.Path)
}

func TestCreateVersionAtExistingVersion(t *testing.T) {
	teardown := setup()
	defer teardown()

	pkg, err := Open(&Config{
		User: "foo",
		Repo: "awesome",
	})
	if err != nil {
		t.Fatal(err)
	}
	v, err := pkg.NewVersion(4)
	assert.Nil(t, v)
	assert.Equal(t, "package foo/awesome: is already at version 4", err.Error())
}
