package main

import (
    "bufio"
    "fmt"
    "os"
)

func canon(x, k int64) int64 {
    r := x % k
    if r == 0 {
        return 0
    }
    if k-r < r {
        return k - r
    }
    return r
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        var k int64
        fmt.Fscan(in, &n, &k)
        freqS := make(map[int64]int)
        for i := 0; i < n; i++ {
            var x int64
            fmt.Fscan(in, &x)
            freqS[canon(x, k)]++
        }
        freqT := make(map[int64]int)
        for i := 0; i < n; i++ {
            var x int64
            fmt.Fscan(in, &x)
            freqT[canon(x, k)]++
        }
        ok := true
        if len(freqS) != len(freqT) {
            ok = false
        } else {
            for key, val := range freqS {
                if freqT[key] != val {
                    ok = false
                    break
                }
            }
        }
        if ok {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
