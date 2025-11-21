package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod = 1000000007

// Precompute factorials up to 2e5
const maxM = 200000

var fact [maxM + 5]int64
var invfact [maxM + 5]int64

func modPow(a, e int64) int64 {
    res := int64(1)
    base := a % mod
    for e > 0 {
        if e&1 == 1 {
            res = res * base % mod
        }
        base = base * base % mod
        e >>= 1
    }
    return res
}

func initFact() {
    fact[0] = 1
    for i := 1; i < len(fact); i++ {
        fact[i] = fact[i-1] * int64(i) % mod
    }
    invfact[len(fact)-1] = modPow(fact[len(fact)-1], mod-2)
    for i := len(fact) - 1; i >= 1; i-- {
        invfact[i-1] = invfact[i] * int64(i) % mod
    }
}

func comb(n, k int) int64 {
    if k < 0 || k > n {
        return 0
    }
    return fact[n] * invfact[k] % mod * invfact[n-k] % mod
}

func main() {
    initFact()

    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        ans := comb(m+n-2, n-1)
        fmt.Fprintln(out, ans)
    }
}

