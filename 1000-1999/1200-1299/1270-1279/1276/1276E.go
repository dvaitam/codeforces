package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var a [4]int64
    var b [4]int64
    for i := 0; i < 4; i++ {
        fmt.Fscan(in, &a[i])
    }
    for i := 0; i < 4; i++ {
        fmt.Fscan(in, &b[i])
    }

    // TODO: implement the algorithm for transforming the set of stones.
    // The solution requires finding a sequence of operations to move from
    // the configuration a to the configuration b using moves of the form
    // choose two stones at x and y and move one from x to 2*y-x.
    // This placeholder always outputs -1.

    fmt.Fprintln(out, -1)
}

