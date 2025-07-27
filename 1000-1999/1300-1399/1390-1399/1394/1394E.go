package main

import (
    "bufio"
    "fmt"
    "os"
)

// Placeholder solution for problem E (maximum number of folds).
// The actual algorithm is non-trivial and not implemented yet.
// This program reads the input and outputs zeros for each prefix
// so that the source compiles successfully.
func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    for i := 0; i < n; i++ {
        var x int
        fmt.Fscan(in, &x)
    }
    for i := 0; i < n; i++ {
        if i > 0 {
            fmt.Fprint(out, " ")
        }
        fmt.Fprint(out, 0)
    }
    fmt.Fprintln(out)
}

