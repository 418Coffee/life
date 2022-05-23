# life

![Tests](https://github.com/418Coffee/life/actions/workflows/test.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/418Coffee/life.svg)](https://pkg.go.dev/github.com/418Coffee/life)
[![Go Report Card](https://goreportcard.com/badge/github.com/418Coffee/life)](https://goreportcard.com/report/github.com/418Coffee/life)

## Table of contents

- [Usage](#usage)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Usage

Install life:

```cmd
go install github.com/418Coffee/life
```

Play around with the cli:

```cmd
life
...
Usage of life.exe [options] width height
options:
  -file string
        load initial state from .rle file (mutually exclusive with width height arguments)
  -nowrap
        don't wrap field toroidally
  -seed int
        seed for initial state (default 1653324678377310)
  -ticks uint
        amount of generation to run (default 100)
```

Or from your own program:

```go
package main

import (
    "fmt"
    "math/rand"
    "time"

    "github.com/418Coffee/life"
)

func main() {
  rand.Seed(0)
  l := life.NewGame(40, 15)
  for i := 0; i < 30; i++ {
    l.Tick()
    fmt.Print("\x1bc", l)
    time.Sleep(time.Second / 30)
   }
}
```

## [Documentation](https://pkg.go.dev/github.com/418Coffee/life)

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)
