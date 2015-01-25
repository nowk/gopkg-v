package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

type ArgError struct {
	msg string
}

func (a ArgError) Error() string {
	return a.msg
}

var (
	ErrInvalidPackagePath    = &ArgError{"invalid package path"}
	ErrInvalidArgumentLength = &ArgError{"invalid arguments length"}
)

type Config struct {
	User string
	Name string

	// Version to create new link at
	Version int
}

func ParseArgs(args []string) (*Config, error) {
	c := &Config{}

	err := parsePackagePath(args[1], c)
	if err != nil {
		return nil, err
	}

	err = parseFlags(args[2:], c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseFlags(args []string, conf *Config) error {
	l := len(args)
	if l < 1 {
		return ErrInvalidArgumentLength
	}

	if nv := args[0]; nv != "--new-version" {
		return &ArgError{fmt.Sprintf("flag provided but not defined: %s", nv)}
	}

	if l == 2 {
		n, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}

		conf.Version = int(n)
	}

	return nil
}

func parsePackagePath(str string, conf *Config) error {
	if str == "" {
		return ErrInvalidPackagePath
	}

	split := strings.Split(str, "/")
	if len(split) != 2 {
		return ErrInvalidPackagePath
	}

	conf.User = split[0]
	conf.Name = split[1]
	if conf.User == "" || conf.Name == "" {
		return ErrInvalidPackagePath
	}

	return nil
}
