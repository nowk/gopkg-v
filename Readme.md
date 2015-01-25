# gopkg-v

[![Build Status](https://travis-ci.org/nowk/gopkg-v.svg?branch=master)](https://travis-ci.org/nowk/gopkg-v)
[![Version](https://img.shields.io/badge/work-in--progress-orange.svg?style=flat)](https://github.com/nowk/gopkg-v)

<!-- [![GoDoc](https://godoc.org/gopkg.in/nowk/gopkg-v.v0?status.svg)](http://godoc.org/gopkg.in/nowk/gopkg-v.v0) -->

gopkg.in version utility to help create and manage versions in relation to the working copy.

---

__The Basic Problem__

Working directory looks something like this:

    github.com/foo/bar
    github.com/foo/bar/baz

And you need to reference a subpackage, but at version.

    package bar

    import "gopkg.in/foo/bar.v1/baz"

To do this you need to create a gopkg.in version that is symlinked to your working copy.

    gopkg.in/foo/bar.v1 => github.com/foo/bar

---

## Install

    go install gopkg.in/nowk/gopkg-v.v0

## Example

Create a new version at the `current version + 1`

    gopkg-v nowk/gopkg-v --new-version
    // gopkg.in/nowk/gopkg-v.v1 -> github.com/nowk/gopkg-v

again...

    gopkg-v nowk/gopkg-v --new-version
    // gopkg.in/nowk/gopkg-v.v2 -> github.com/nowk/gopkg-v

* If there is no existing versions, it will start at `v0`
* It will remove the previous version link, if it is a link. Physical directories will be ignored.

---

Create a new version at a specific version

    gopkg-v nowk/gopkg-v --new-version 3
    // gopkg.in/nowk/gopkg-v.v3 -> github.com/nowk/gopkg-v

* If `--new-version 0` it will attempt to create at `v0` if there are no existing versions. Else it will create a link at the `current version + 1`.


## License

MIT

