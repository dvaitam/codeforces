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
        var n, x int
        fmt.Fscan(in, &n, &x)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }

        pref := make([]int, n+1)
        for i := 1; i <= n; i++ {
            pref[i] = pref[i-1] + a[i-1]
        }

        maxSum := make([]int, n+1)
        for length := 1; length <= n; length++ {
            best := -1 << 60
            for i := 0; i+length <= n; i++ {
                s := pref[i+length] - pref[i]
                if s > best {
                    best = s
                }
            }
            maxSum[length] = best
        }

        for k := 0; k <= n; k++ {
            ans := 0
            for length := 0; length <= n; length++ {
                add := k
                if length < k {
                    add = length
                }
                val := maxSum[length] + add*x
                if val > ans {
                    ans = val
                }
            }
            if k > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, ans)
        }
        fmt.Fprintln(out)
    }
}
