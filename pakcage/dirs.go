package pakcage

import (
	"os"
	"regexp"
)

type dirnames []string

func (d dirnames) Len() int {
	return len(d)
}

func (d dirnames) Last() (string, bool) {
	if n := d.Len(); n > 0 {
		return d[n-1], true
	}

	return "", false
}

// matchDirs matches child (one level deep) directories to a given regex
func matchDirs(f *os.File, reg *regexp.Regexp) (dirnames, error) {
	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	var d dirnames
	for _, v := range names {
		if reg.MatchString(v) {
			d = append(d, v)
		}
	}

	return d, nil
}
