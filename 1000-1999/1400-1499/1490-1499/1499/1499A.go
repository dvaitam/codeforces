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

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, k1, k2 int
        fmt.Fscan(in, &n, &k1, &k2)
        var w, b int
        fmt.Fscan(in, &w, &b)

        whiteCells := k1 + k2
        maxWhiteDominoes := whiteCells / 2
        blackCells := 2*n - whiteCells
        maxBlackDominoes := blackCells / 2

        if w <= maxWhiteDominoes && b <= maxBlackDominoes {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}

