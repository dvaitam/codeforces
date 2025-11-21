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
        var n int
        fmt.Fscan(in, &n)
        sum := int64(0)
        odd := 0
        maxVal := int64(0)
        for i := 0; i < n; i++ {
            var x int64
            fmt.Fscan(in, &x)
            sum += x
            if x&1 == 1 {
                odd++
            }
            if x > maxVal {
                maxVal = x
            }
        }
        if odd == 0 || odd == n {
            fmt.Fprintln(out, maxVal)
        } else {
            ans := sum - int64(odd-1)
            fmt.Fprintln(out, ans)
        }
    }
}
