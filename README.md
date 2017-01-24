[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/txgruppi/retry-go)
![Codeship](https://img.shields.io/codeship/3e5c23f0-c2d6-0133-421c-025eedb952b8.svg?style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/github/txgruppi/retry-go.svg?style=flat-square)](https://codecov.io/github/txgruppi/retry-go)
[![Go Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg?style=flat-square)](https://goreportcard.com/report/github.com/txgruppi/retry-go)

# Retry

Retry a function execution with specific intervals with panic recovery

Make sure to read the docs to understand how this package works and what do
expected from it.

## Installation

```
go get -u github.com/txgruppi/retry-go
```

## Example

```go
package main

import (
  "log"
  "time"

  "github.com/txgruppi/retry-go"
)

func main() {
  toTry := func(attempt, limit int) error {
    // Do you stuff and return an error if there is any
    return nil
  }

  err := retry.WithFixedInterval(1 * time.Second, 5, toTry)
  if err != nil {
    log.Fatal(err) // It should log the errors if there were any
  }
}
```

## Tests

```
go get -u -t github.com/txgruppi/retry-go
cd $GOPATH/src/github.com/txgruppi/retry-go
go test ./...
```

## License

MIT
