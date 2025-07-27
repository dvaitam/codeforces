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
        var n, k int64
        fmt.Fscan(in, &n, &k)
        denom := int64(1) + k + k*k + k*k*k
        n1 := n / denom
        n2 := k * n1
        n3 := k * n2
        n4 := k * n3
        fmt.Fprintf(out, "%d %d %d %d\n", n1, n2, n3, n4)
    }
}
