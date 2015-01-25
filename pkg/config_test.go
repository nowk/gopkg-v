package pkg

import (
	"gopkg.in/nowk/assert.v2"
	"testing"
)

func TestParseValidConfig(t *testing.T) {
	args := []string{"/path/to/bin", "nowk/gopkg-v", "--gopkg-version", "1"}
	config, err := ParseArgs(args)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &Config{
		User:    "nowk",
		Repo:    "gopkg-v",
		Version: 1,
	}, config)
}

func TestInvalidPackagePath(t *testing.T) {
	for _, v := range [][]string{
		{"/path", ""},
		{"/path", "nowk"},
		{"/path", "nowk/"},
		{"/path", "/"},
		{"/path", "/gopkg-v"},
	} {
		_, err := ParseArgs(v)
		assert.Equal(t, ErrInvalidPackagePath, err)
	}
}

func TestArgErrors(t *testing.T) {
	for _, v := range []struct {
		Args  []string
		Error string
	}{
		{
			[]string{"/path/to/bin", "nowk/gopkg-v"},
			"invalid arguments length",
		},
		{
			[]string{"/path/to/bin", "nowk/gopkg-v", "--gopkg-versions", "1"},
			"flag provided but not defined: -gopkg-versions",
		},
	} {
		_, err := ParseArgs(v.Args)
		assert.Equal(t, v.Error, err.Error())
	}
}
