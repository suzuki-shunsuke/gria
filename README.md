# Gria

CLI tool for golang's test code scaffolding.

## Overview

gria is CLI tool for golang's test code scaffolding.
Parse files and packages with [go/ast](https://golang.org/pkg/go/ast) and generate test functions skeltons.

## Getting Started

First, [install gria]().

Create a sample code `foo.go`.

```golang
package foo

func Foo() string {
	return "foo"
}
```

Then run `gria gen`.

```
$ gria gen foo.go
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

Then run `gria gen` again.

```
$ gria gen foo.go
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

Notice the existing code is not overwritten.

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

## Configuration

```yaml
---
ignore_unxposed: false  # ignore unexposed functions
func_name: Test{{.StructName}}{{.Name}}  # test function name
includes:  # white list
excludes:  # black list
  foo: # ignore package
  zoo/bar.go: # ignore file
  config/config.go:
  - GetName # ignore function
  - User.GetAge # ignore method
  - Foo # ignore struct
```

## License

[MIT](LICENSE)
