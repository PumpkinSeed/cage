# Cage

Simple library to catch stdout/stderr in Go.

#### Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/PumpkinSeed/cage"
)

func main() {
    c := cage.Start()
    
    fmt.Println("test")
    fmt.Println("test2")
    fmt.Fprintln(os.Stderr, "stderr error")
    
    cage.Stop(c)
    fmt.Println(c.Data)
    // [test, test2, stderr error]
}
```
