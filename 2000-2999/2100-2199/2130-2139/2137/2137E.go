package main

import (
    "bufio"
    "fmt"
    "os"
)

const MOD int64 = 998244353

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        var k int64
        fmt.Fscan(in, &n, &k)
        a := make([]int64, n)
        sum := int64(0)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            sum += a[i]
        }
        fmt.Fprintln(out, sum%MOD)
    }
}
