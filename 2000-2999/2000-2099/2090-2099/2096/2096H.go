package main

import (
    "bufio"
    "fmt"
    "os"
)

const MOD int64 = 998244353

// fast exponentiation modulo MOD.
func modPow(a, e int64) int64 {
    res := int64(1)
    for e > 0 {
        if e&1 == 1 {
            res = res * a % MOD
        }
        a = a * a % MOD
        e >>= 1
    }
    return res
}

// Walsh-Hadamard transform for XOR convolution.
func fwht(a []int64, invert bool) {
    n := len(a)
    for len := 1; 2*len <= n; len <<= 1 {
        for i := 0; i < n; i += 2 * len {
            for j := 0; j < len; j++ {
                u := a[i+j]
                v := a[i+j+len]
                a[i+j] = (u + v) % MOD
                a[i+j+len] = (u - v + MOD) % MOD
            }
        }
    }
    if invert {
        invN := modPow(int64(n), MOD-2)
        for i := 0; i < n; i++ {
            a[i] = a[i] * invN % MOD
        }
    }
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }

    for ; T > 0; T-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        size := 1 << m

        // compress identical intervals
        type key struct{ l, r int }
        mp := make(map[key]int64)
        intervals := make([]key, 0)
        for i := 0; i < n; i++ {
            var l, r int
            fmt.Fscan(in, &l, &r)
            k := key{l, r}
            if _, ok := mp[k]; !ok {
                intervals = append(intervals, k)
            }
            mp[k]++
        }

        // product in Walsh domain
        prod := make([]int64, size)
        for i := 0; i < size; i++ {
            prod[i] = 1
        }

        // precompute powers of two mod
        pow2 := make([]int64, size)
        pow2[0] = 1
        for i := 1; i < size; i++ {
            pow2[i] = pow2[i-1] * 2 % MOD
        }

        for _, k := range intervals {
            arr := make([]int64, size)
            for i := k.l; i <= k.r; i++ {
                arr[i] = 1
            }
            fwht(arr, false)
            cnt := mp[k]
            for i := 0; i < size; i++ {
                arr[i] = modPow(arr[i]%MOD, cnt)
                prod[i] = prod[i] * arr[i] % MOD
            }
        }

        fwht(prod, true)

        var h uint64
        for x := 0; x < size; x++ {
            g := prod[x] * pow2[x] % MOD
            h ^= uint64(g)
        }
        fmt.Fprintln(out, h)
    }
}

