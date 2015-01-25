package pkg

import (
	"errors"
	"flag"
	"strings"
)

var (
	ErrInvalidPackagePath    = errors.New("invalid package path")
	ErrInvalidArgumentLength = errors.New("invalid arguments length")
)

type Config struct {
	User    string
	Repo    string
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
	if len(args) < 1 {
		return ErrInvalidArgumentLength
	}

	vflag := flag.NewFlagSet("-", flag.ContinueOnError)
	vflag.IntVar(&conf.Version, "gopkg-version", -1, "")
	return vflag.Parse(args)
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
	conf.Repo = split[1]
	if conf.User == "" || conf.Repo == "" {
		return ErrInvalidPackagePath
	}

	return nil
}
