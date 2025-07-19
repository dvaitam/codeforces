package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod = int64(1e9 + 7)

func modExp(x, y int64) int64 {
    res := int64(1)
    for y > 0 {
        if y&1 == 1 {
            res = res * x % mod
        }
        x = x * x % mod
        y >>= 1
    }
    return res
}

func pd(c byte) int {
    if c >= 'A' && c <= 'Z' {
        return int(c - 'A' + 1)
    }
    return int(c - 'a' + 27)
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var s string
    var q int
    if _, err := fmt.Fscan(reader, &s); err != nil {
        return
    }
    n := len(s)
    fmt.Fscan(reader, &q)

    // count letters
    cnt := make([]int, 53)
    for i := 0; i < n; i++ {
        cnt[pd(s[i])]++
    }

    // factorial and inverse factorial
    jc := make([]int64, n+1)
    inv := make([]int64, n+1)
    jc[0] = 1
    for i := 1; i <= n; i++ {
        jc[i] = jc[i-1] * int64(i) % mod
    }
    inv[n] = modExp(jc[n], mod-2)
    for i := n; i > 0; i-- {
        inv[i-1] = inv[i] * int64(i) % mod
    }

    half := n / 2
    // base multiplier
    g := jc[half] * jc[half] % mod
    for i := 1; i <= 52; i++ {
        if cnt[i] > 0 {
            g = g * inv[cnt[i]] % mod
        }
    }

    // subset-sum dp
    p := make([]int64, n+1)
    p[0] = 1
    for i := 1; i <= 52; i++ {
        c := cnt[i]
        if c == 0 {
            continue
        }
        for j := n; j >= c; j-- {
            p[j] = (p[j] + p[j-c]) % mod
        }
    }

    // prepare results
    var sres [53][53]int64
    p2 := make([]int64, n+1)
    // for single letter i,i
    for i := 1; i <= 52; i++ {
        ci := cnt[i]
        if ci == 0 || ci > half {
            continue
        }
        // p2 = p without type i
        copy(p2, p)
        for j := ci; j <= n; j++ {
            p2[j] = (p2[j] - p2[j-ci] + mod) % mod
        }
        idx := half - ci
        if idx >= 0 {
            sres[i][i] = p2[idx] * g % mod * 2 % mod
        }
        // pairs i < j
        for j := i + 1; j <= 52; j++ {
            cj := cnt[j]
            if cj == 0 || ci+cj > half {
                continue
            }
            // compute alternating sum
            w := half - ci - cj
            var ct int64
            sign := int64(1)
            for ww := w; ww >= 0; ww -= cj {
                ct = (ct + p2[ww]*sign) % mod
                sign = -sign
            }
            sres[i][j] = (ct%mod+mod)%mod * g % mod * 2 % mod
        }
    }

    // answer queries
    for k := 0; k < q; k++ {
        var x, y int
        fmt.Fscan(reader, &x, &y)
        xi := pd(s[x-1])
        yi := pd(s[y-1])
        if xi > yi {
            xi, yi = yi, xi
        }
        res := sres[xi][yi]
        fmt.Fprintln(writer, res)
    }
}
