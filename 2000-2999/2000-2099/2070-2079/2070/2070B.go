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

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        var x int
        var k int64
        fmt.Fscan(in, &n, &x, &k)
        var s string
        fmt.Fscan(in, &s)

        pref := make([]int, n+1)
        first := make(map[int]int, 2*n+1)
        for i := 0; i < n; i++ {
            delta := -1
            if s[i] == 'R' {
                delta = 1
            }
            pref[i+1] = pref[i] + delta
            if _, ok := first[pref[i+1]]; !ok {
                first[pref[i+1]] = i + 1
            }
        }

        firstZeroFrom := func(start int) (int, bool) {
            target := -start
            idx, ok := first[target]
            if !ok {
                return 0, false
            }
            return idx, true
        }

        t1, ok := firstZeroFrom(x)
        if !ok || int64(t1) > k {
            fmt.Fprintln(out, 0)
            continue
        }
        count := int64(1)
        k -= int64(t1)

        t0, ok := firstZeroFrom(0)
        if !ok {
            fmt.Fprintln(out, count)
            continue
        }
        count += k / int64(t0)
        fmt.Fprintln(out, count)
    }
}
