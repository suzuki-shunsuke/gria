# Gria

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
