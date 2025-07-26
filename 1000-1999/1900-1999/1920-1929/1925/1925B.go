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
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var x, n int64
        fmt.Fscan(in, &x, &n)
        maxDiv := int64(1)
        limit := x / n
        for i := int64(1); i*i <= x; i++ {
            if x%i == 0 {
                d1 := i
                d2 := x / i
                if d1 <= limit && d1 > maxDiv {
                    maxDiv = d1
                }
                if d2 <= limit && d2 > maxDiv {
                    maxDiv = d2
                }
            }
        }
        fmt.Fprintln(out, maxDiv)
    }
}

