package main

import (
    "bufio"
    "fmt"
    "os"
)

const MOD int64 = 1000000007

// Precompute factorials up to maxA and their inverses.
var fact, invFact []int64

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

// Combination C(n, r) mod MOD for r < MOD (here r <= 1e5).
// If the interval [n-r+1, n] contains a multiple of MOD, the result is 0.
func combLarge(n int64, r int) int64 {
    if r < 0 || int64(r) > n {
        return 0
    }
    if r == 0 || r == 1 {
        if r == 0 {
            return 1
        }
        return n % MOD
    }

    // Check if a multiple of MOD lies in [n-r+1, n].
    // Count of multiples up to x is x / MOD.
    multiples := n/MOD - (n-int64(r))/MOD
    if multiples > 0 {
        return 0
    }

    num := int64(1)
    for i := 0; i < r; i++ {
        num = num * ((n - int64(i)) % MOD) % MOD
    }
    num = num * invFact[r] % MOD
    return num
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)

    cases := make([][3]int64, T)
    maxA := 0
    for i := 0; i < T; i++ {
        var a, b, k int64
        fmt.Fscan(in, &a, &b, &k)
        cases[i] = [3]int64{a, b, k}
        if int(a) > maxA {
            maxA = int(a)
        }
    }

    // Precompute factorials up to maxA.
    fact = make([]int64, maxA+1)
    invFact = make([]int64, maxA+1)
    fact[0] = 1
    for i := 1; i <= maxA; i++ {
        fact[i] = fact[i-1] * int64(i) % MOD
    }
    invFact[maxA] = modPow(fact[maxA], MOD-2)
    for i := maxA - 1; i >= 0; i-- {
        invFact[i] = invFact[i+1] * int64(i+1) % MOD
    }

    for _, tc := range cases {
        a := tc[0]
        b := tc[1]
        k := tc[2]

        // Minimal rows needed so that some color must appear in a rows of a column.
        n := k*(a-1) + 1

        // Choose columns so that pigeonhole forces at least b columns with same color on the same set of a rows.
        comb := combLarge(n, int(a))
        m := (k % MOD) * ((b - 1) % MOD) % MOD
        m = (m * comb) % MOD
        m = (m + 1) % MOD

        fmt.Fprintf(out, "%d %d\n", n%MOD, m)
    }
}
