package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod int64 = 1_000_000_007
const size = 1 << 10

var inv2 int64
var inv10000 int64

func modPow(a, b int64) int64 {
    res := int64(1)
    for b > 0 {
        if b&1 == 1 {
            res = res * a % mod
        }
        a = a * a % mod
        b >>= 1
    }
    return res
}

func init() {
    inv2 = (mod + 1) / 2
    inv10000 = modPow(10000, mod-2)
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        p := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &p[i])
        }

        prod := make([]int64, size)
        used := make([]bool, size)
        for i := 0; i < size; i++ {
            prod[i] = 1
        }
        values := make([]int, 0, size)

        for i := 0; i < n; i++ {
            v := a[i]
            if !used[v] {
                used[v] = true
                values = append(values, v)
            }
            qi := int64(p[i]) * inv10000 % mod
            term := (1 - (2*qi%mod) + mod) % mod
            prod[v] = prod[v] * term % mod
        }

        dp := make([]int64, size)
        ndp := make([]int64, size)
        dp[0] = 1

        for _, v := range values {
            odd := (1 - prod[v] + mod) % mod
            odd = odd * inv2 % mod
            if odd == 0 {
                continue
            }
            stay := (1 - odd + mod) % mod
            for i := 0; i < size; i++ {
                val := (dp[i]*stay + dp[i^v]*odd) % mod
                ndp[i] = val
            }
            dp, ndp = ndp, dp
        }

        ans := int64(0)
        for x := 0; x < size; x++ {
            val := int64(x)
            ans = (ans + dp[x]*val%mod*val) % mod
        }
        fmt.Fprintln(out, ans)
    }
}
