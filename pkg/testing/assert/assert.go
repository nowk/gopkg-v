package assert

import (
	"fmt"
	"gopkg.in/nowk/assert.v2"
	"os"
	"testing"
)

var (
	Equal = assert.Equal
	True  = assert.True
	False = assert.False
	Nil   = assert.Nil
)

func Symlink(t *testing.T, exp, got string) {
	f, err := os.Readlink(got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exp, f)
}

func Unlink(t *testing.T, path string) {
	_, err := os.Open(path)
	assert.Equal(t, fmt.Sprintf("open %s: no such file or directory", path),
		err.Error())
}
