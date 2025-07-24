package main

import (
    "bufio"
    "fmt"
    "os"
)

// This is a placeholder implementation for problem H2 from contest 1781.
// The full combinatorial solution has not been translated to Go yet.
// The program reads the input format as specified and prints zero for
// each test case so that it compiles and can serve as a starting point
// for future development.

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var h, w, k int
        fmt.Fscan(in, &h, &w, &k)
        for i := 0; i < k; i++ {
            var r, c int
            fmt.Fscan(in, &r, &c)
        }
        // TODO: implement the actual logic to count distinct signals.
        fmt.Fprintln(out, 0)
    }
}

