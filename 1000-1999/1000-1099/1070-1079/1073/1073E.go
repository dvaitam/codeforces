package main

import (
    "bufio"
    "fmt"
    "math/bits"
    "os"
)

const (
    MOD = 998244353
    N   = 25
    K   = 10
)

var (
    p      [N]int
    f      [N][1 << K]int
    g      [N][1 << K]int
    kParam int
)

func add(a, b int) int {
    a += b
    if a >= MOD {
        a -= MOD
    }
    return a
}

func sub(a, b int) int {
    a -= b
    if a < 0 {
        a += MOD
    }
    return a
}

func mul(a, b int) int {
    return int((int64(a) * int64(b)) % MOD)
}

func initDP() {
    p[0] = 1
    for i := 1; i < N; i++ {
        p[i] = mul(p[i-1], 10)
    }
    g[0][0] = 1
    for i := 1; i < N; i++ {
        for j := 0; j < 10; j++ {
            for mask := 0; mask < (1 << K); mask++ {
                newMask := mask | (1 << j)
                g[i][newMask] = add(g[i][newMask], g[i-1][mask])
                foo := mul(j, p[i-1])
                foo = mul(foo, g[i-1][mask])
                foo = add(foo, f[i-1][mask])
                f[i][newMask] = add(f[i][newMask], foo)
            }
        }
    }
}

// get returns the sum of all numbers x < val with at most kParam distinct digits
func get(val int64) int {
    if val <= 0 {
        return 0
    }
    var v []int
    for val > 0 {
        v = append(v, int(val%10))
        val /= 10
    }
    res := 0
    curMask := 0
    tot := 0
    sz := len(v)

    // sum for numbers with fewer digits
    for i := 1; i < sz; i++ {
        for j := 1; j <= 9; j++ {
            for mask := 0; mask < (1 << K); mask++ {
                newMask := mask | (1 << j)
                if bits.OnesCount(uint(newMask)) > kParam {
                    continue
                }
                foo := mul(j, p[i-1])
                foo = mul(foo, g[i-1][mask])
                foo = add(foo, f[i-1][mask])
                res = add(res, foo)
            }
        }
    }

    // sum for numbers with same number of digits and less than val
    for i := sz; i > 0; i-- {
        u := v[i-1]
        start := 0
        if i == sz {
            start = 1
        }
        for j := start; j < u; j++ {
            newMask := curMask | (1 << j)
            newTot := add(tot, mul(j, p[i-1]))
            for mask := 0; mask < (1 << K); mask++ {
                fooMask := newMask | mask
                if bits.OnesCount(uint(fooMask)) > kParam {
                    continue
                }
                bar := mul(newTot, g[i-1][mask])
                bar = add(bar, f[i-1][mask])
                res = add(res, bar)
            }
        }
        // include digit u
        curMask |= 1 << u
        tot = add(tot, mul(u, p[i-1]))
    }
    return res
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    var l, r int64
    fmt.Fscan(reader, &l, &r, &kParam)
    initDP()
    ans := sub(get(r+1), get(l))
    fmt.Println(ans)
}
