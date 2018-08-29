# go-bringo

Golang wrapper for bringo247.ru API delivery

## Install

Install the package with:

```bash
go get github.com/xMlex/go-bringo
```

Import it with:

```go
import "github.com/xMlex/go-bringo"
```

and use `go-bringo` as the package name inside the code.

## Example

```go
package main

import (
	"github.com/xMlex/go-bringo"
	"log"
	"strconv"
)

func main() {
	api := bringo.New()
    api.Init("login", "password", false) // false - is Production api url

    if err := api.Login(); err == nil {
		log.Fatalln(err)
	}

    // methods
    // api.Calculate(&bringo.Delivery{})
    // api.Create(&bringo.Delivery{})
    // api.Info(1)
    // api.Cancel(1)
    // api.AccountInfo
}
```

and see `bringo_test.go` for more examples