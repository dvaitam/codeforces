package main

import (
    "bufio"
    "fmt"
    "os"
)

const MOD int = 998244353
const ROOT int = 3

func modAdd(a, b int) int { a += b; if a >= MOD { a -= MOD }; return a }
func modSub(a, b int) int { a -= b; if a < 0 { a += MOD }; return a }
func modMul(a, b int) int { return int(int64(a) * int64(b) % int64(MOD)) }
func modPow(a, e int) int { res := 1; for e > 0 { if e&1 == 1 { res = modMul(res, a) }; a = modMul(a, a); e >>= 1 }; return res }
func modInv(a int) int { return modPow(a, MOD-2) }

func ntt(a []int, invert bool) {
    n := len(a)
    j := 0
    for i := 1; i < n; i++ {
        bit := n >> 1
        for ; j&bit != 0; bit >>= 1 { j ^= bit }
        j |= bit
        if i < j {
            a[i], a[j] = a[j], a[i]
        }
    }
    for length := 2; length <= n; length <<= 1 {
        wlen := modPow(ROOT, (MOD-1)/length)
        if invert { wlen = modInv(wlen) }
        half := length >> 1
        for i := 0; i < n; i += length {
            w := 1
            for j := 0; j < half; j++ {
                u := a[i+j]
                v := modMul(a[i+j+half], w)
                a[i+j] = modAdd(u, v)
                a[i+j+half] = modSub(u, v)
                w = modMul(w, wlen)
            }
        }
    }
    if invert {
        invN := modInv(n)
        for i := 0; i < n; i++ {
            a[i] = modMul(a[i], invN)
        }
    }
}

func polyMul(a, b []int) []int {
    n := len(a) + len(b) - 1
    sz := 1
    for sz < n { sz <<= 1 }
    fa := make([]int, sz)
    fb := make([]int, sz)
    copy(fa, a)
    copy(fb, b)
    ntt(fa, false)
    ntt(fb, false)
    for i := 0; i < sz; i++ {
        fa[i] = modMul(fa[i], fb[i])
    }
    ntt(fa, true)
    return fa[:n]
}

func solve(n int) []int {
    k := (n - 1) / 2
    m := (n + 1) / 2
    fact := make([]int, n+1)
    fact[0] = 1
    for i := 1; i <= n; i++ {
        fact[i] = modMul(fact[i-1], i)
    }
    invFact := make([]int, n+1)
    invFact[n] = modInv(fact[n])
    for i := n; i >= 1; i-- {
        invFact[i-1] = modMul(invFact[i], i)
    }
    inv := make([]int, n+1)
    for i := 1; i <= n; i++ {
        inv[i] = modInv(i)
    }
    F := make([]int, n)
    pref := make([]int, n)
    F[0] = 1
    pref[0] = 1
    for x := 1; x < n; x++ {
        l := x - k
        if l < 0 {
            l = 0
        }
        sum := pref[x-1]
        if l-1 >= 0 {
            sum = modSub(sum, pref[l-1])
        }
        F[x] = modMul(sum, inv[x])
        pref[x] = modAdd(pref[x-1], F[x])
    }
    f := make([]int, n)
    for i := 0; i < n; i++ {
        f[i] = modMul(F[i], fact[i])
    }
    A := make([]int, n)
    for t := m - 1; t <= n-2; t++ {
        A[t] = modMul(F[t], fact[n-t-2])
    }
    B := make([]int, n)
    for v := 0; v < n; v++ {
        B[v] = invFact[v]
    }
    D := polyMul(A, B)
    ans := make([]int, n)
    ans[0] = f[n-1]
    for i := 2; i <= n; i++ {
        if i > m {
            ans[i-1] = 0
            continue
        }
        u := n - i
        val := 0
        if u >= 0 && u < len(D) {
            val = D[u]
        }
        term := modMul(modMul(i-1, fact[n-i]), val)
        ans[i-1] = term
    }
    for i := m + 1; i <= n; i++ {
        if i-1 < n {
            ans[i-1] = 0
        }
    }
    return ans
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var n int
    fmt.Fscan(reader, &n)
    res := solve(n)
    for i, v := range res {
        if i > 0 {
            fmt.Fprint(writer, " ")
        }
        fmt.Fprint(writer, v)
    }
    fmt.Fprintln(writer)
}

