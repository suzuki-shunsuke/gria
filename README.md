# Gria

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/suzuki-shunsuke/gria)
[![Build Status](https://travis-ci.org/suzuki-shunsuke/gria.svg?branch=master)](https://travis-ci.org/suzuki-shunsuke/gria)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/gria/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/gria)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/gria)](https://goreportcard.com/report/github.com/suzuki-shunsuke/gria)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/gria.svg)](https://github.com/suzuki-shunsuke/gria)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/gria.svg)](https://github.com/suzuki-shunsuke/gria/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/gria/master/LICENSE)

CLI tool for golang's test code scaffolding.

## Overview

gria is CLI tool for golang's test code scaffolding.
Parse packages with [go/ast](https://golang.org/pkg/go/ast) and generate test functions skeltons.

## Getting Started

First, [install gria](#install).

Create a sample code `foo.go`.

```golang
package foo

func Foo() string {
	return "foo"
}
```

Then run `gria`.

```
$ gria .
create a test file: foo_test.go
```

Open the `foo_test.go`.

```golang
package foo

import (
	"testing"
)

func TestFoo(t *testing.T) {
}
```

Then add a struct and method.

```golang
package foo

func Foo() string {
	return "foo"
}

type User struct {
	name string
}

func (user User) GetName() string {
	return user.name
}
```

Then run `gria` again.

```
$ gria .
add a test skelton code to a test file: foo_test.go
```

Open the `foo_test.go`.

```golang
package foo

import (
	"testing"
)

func TestFoo(t *testing.T) {
}

func TestUserGetName(t *testing.T) {
}
```

Note the existing code is not overwritten.

## Install

```
$ go get github.com/suzuki-shunsuke/gria
```

or Download a binary from the [release page](https://github.com/suzuki-shunsuke/gria/releases).

Check whether gria is installed.

```
$ gria -v
gria version 0.1.0
```

## Usage

Show help.

```
$ gria --h
$ gria --help
$ gria h
$ gria help
```

Show version.

```
$ gria -v
$ gria --version
```

Generate tests.

```
$ gria <package_directory_path> [<package_directory_path> ...]
```

If you want to pass packages recursively, use `go list`.

```
$ go list -f "{{.Dir}}" ./... | xargs gria
```

## License

[MIT](LICENSE)
