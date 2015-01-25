package pkg

import (
	"gopkg.in/nowk/gopkg-v.v0/pkg/testing/assert"
	"testing"
)

func TestParseValidConfig(t *testing.T) {
	args := []string{"/path/to/bin", "nowk/gopkg-v", "--new-version", "1"}
	config, err := ParseArgs(args)
	assert.Nil(t, err)
	assert.Equal(t, &Config{
		User:    "nowk",
		Name:    "gopkg-v",
		Version: 1,
	}, config)
}

func TestOnlyNewVersionFlag(t *testing.T) {
	args := []string{"/path/to/bin", "nowk/gopkg-v", "--new-version"}
	config, err := ParseArgs(args)
	assert.Nil(t, err)
	assert.Equal(t, &Config{
		User:    "nowk",
		Name:    "gopkg-v",
		Version: 0,
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
			[]string{"/path/to/bin", "nowk/gopkg-v", "--new-versions", "1"},
			"flag provided but not defined: --new-versions",
		},
	} {
		_, err := ParseArgs(v.Args)
		assert.Equal(t, v.Error, err.Error())
	}
}
