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
        arr := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &arr[i])
        }
        ans := (1 << 30) - 1
        for i, v := range arr {
            if v != i {
                ans &= v
            }
        }
        fmt.Fprintln(out, ans)
    }
}
