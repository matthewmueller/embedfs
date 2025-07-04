# embedfs

[![Go Reference](https://pkg.go.dev/badge/github.com/matthewmueller/embedfs.svg)](https://pkg.go.dev/github.com/matthewmueller/embedfs)

Easily switch between live files in development and embedded files in production.

## Features

- Works with `go run` and `go build`.
- Tiny wrapper on top of `embed.FS`
- No need to change your code or environment.
- Minimal dependencies and easy to use.

## Example

```go
import (
  "embed"
  "github.com/matthewmueller/embedfs"
  "fmt"
)

//go:embed public/**
var publicfs embed.FS

func main() {
  // In development, when using `go run`, this code reads files in `./public`.
  // In production, when you start a binary built with `go build`, this code
  // uses the embedded filesystem (`publicfs` in this case).
  fsys, err := embedfs.Load(publicfs, "public")
  if err != nil {
    // ...
  }

  // Serve the files
  data, err := fs.ReadFile(fsys, "public/favicon.ico")
  if err != nil {
    // ...
  }
  fmt.Println(string(data))
}
```

## Development

First, clone the repo:

```sh
git clone https://github.com/matthewmueller/embedfs
cd embedfs
```

Next, install dependencies:

```sh
go mod tidy
```

Finally, try running the tests:

```sh
go test ./...
```

## License

MIT
