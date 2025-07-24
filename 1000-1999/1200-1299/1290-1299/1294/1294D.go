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

    var q, x int
    fmt.Fscan(in, &q, &x)
    cnt := make([]int, x)
    mex := 0
    for i := 0; i < q; i++ {
        var y int
        fmt.Fscan(in, &y)
        cnt[y%x]++
        for cnt[mex%x] > 0 {
            cnt[mex%x]--
            mex++
        }
        fmt.Fprintln(out, mex)
    }
}
