# gopkg-v

[![Build Status](https://travis-ci.org/nowk/gopkg-v.svg?branch=master)](https://travis-ci.org/nowk/gopkg-v)
[![Version](https://img.shields.io/badge/work-in--progress-orange.svg?style=flat)](https://github.com/nowk/gopkg-v)

<!-- [![GoDoc](https://godoc.org/gopkg.in/nowk/gopkg-v.v0?status.svg)](http://godoc.org/gopkg.in/nowk/gopkg-v.v0) -->

gopkg.in version utility

__*Work In Progress*__

---

Sets up symbolic links to your __gopkg.in__ dir to your __github.com__ repo directory. Allowing one to work within a git repo and use the proper import directory structure for gopkg versioning.

Working directory looks something like this:

    github.com/foo/bar
    github.com/foo/bar/baz

Version directory is symlinked:

    gopkg.in/foo/bar.v1 => github.com/foo/bar

Now you can import subpackages within the version context of your parent package.

    package bar

    import "gopkg.in/foo/bar.v1/baz"

---

## Install

    go install gopkg.in/nowk/gopkg-v.v0

## License

MIT

