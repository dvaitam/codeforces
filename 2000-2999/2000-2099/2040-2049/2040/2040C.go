package main

import (
    "bufio"
    "fmt"
    "os"
)

const maxPowVal int64 = 1e18

func buildPow2(limit int) []int64 {
    pow := make([]int64, limit+1)
    pow[0] = 1
    for i := 1; i <= limit; i++ {
        val := pow[i-1] << 1
        if val > maxPowVal {
            val = maxPowVal
        }
        pow[i] = val
    }
    return pow
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)

    maxN := 200000
    pow2 := buildPow2(maxN)

    for ; t > 0; t-- {
        var n int
        var k int64
        fmt.Fscan(in, &n, &k)

        if pow2[n-1] < k {
            fmt.Fprintln(out, -1)
            continue
        }

        left := make([]int, 0)
        kLeft := k

        for x := 1; x <= n-1; x++ {
            count := pow2[n-1-x]
            if kLeft > count {
                kLeft -= count
                continue
            }
            left = append(left, x)
            for y := x + 1; y <= n-1; y++ {
                count = pow2[n-1-y]
                if kLeft > count {
                    kLeft -= count
                    continue
                }
                left = append(left, y)
            }
            break
        }

        used := make([]bool, n+1)
        for _, v := range left {
            used[v] = true
        }

        perm := make([]int, 0, n)
        perm = append(perm, left...)
        perm = append(perm, n)
        for v := n - 1; v >= 1; v-- {
            if !used[v] {
                perm = append(perm, v)
            }
        }

        for i, val := range perm {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, val)
        }
        fmt.Fprintln(out)
    }
}
