package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod int64 = 1000000007

func modPow(a, e int64) int64 {
    res := int64(1)
    for e > 0 {
        if e&1 == 1 {
            res = res * a % mod
        }
        a = a * a % mod
        e >>= 1
    }
    return res
}

func main() {
    // Precompute modular inverses for all possible (l + k).
    const maxLK = 5050
    inv := make([]int64, maxLK+1)
    for i := 1; i <= maxLK; i++ {
        inv[i] = modPow(int64(i), mod-2)
    }

    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n, l, k int
        fmt.Fscan(in, &n, &l, &k)

        // dp[k][i] — probability the current state has k counterfeit keys remaining,
        // next turn belongs to player i.
        dp := make([][]int64, k+1)
        for i := range dp {
            dp[i] = make([]int64, n)
        }
        dp[k][0] = 1

        ans := make([]int64, n)

        for rem := l; rem >= 1; rem-- {
            nextR := make([][]int64, k+1)
            for i := range nextR {
                nextR[i] = make([]int64, n)
            }

            for currK := k; currK >= 0; currK-- {
                state := dp[currK]
                // total probability mass of this state
                var total int64
                for _, v := range state {
                    total = (total + v) % mod
                }
                if total == 0 {
                    continue
                }

                invAll := inv[rem+currK]            // 1 / (rem + currK)
                winCoeff := invAll                   // w = 1 / (rem + currK)
                fakeProb := int64(currK) * invAll % mod

                q := rem / n
                rMod := rem % n
                qMod := int64(q) % mod

                // addMore[j]: sum of state over a window of length rMod ending at position j (circular)
                addMore := make([]int64, n)
                if rMod > 0 {
                    ext := make([]int64, 2*n+1)
                    for i := 0; i < 2*n; i++ {
                        ext[i+1] = (ext[i] + state[i%n]) % mod
                    }
                    for j := 0; j < n; j++ {
                        lIdx := j + n - rMod + 1
                        rIdx := j + n + 1
                        add := ext[rIdx] - ext[lIdx]
                        if add < 0 {
                            add += mod
                        }
                        addMore[j] = add
                    }
                }

                // Real key case: distribute wins and next start.
                for j := 0; j < n; j++ {
                    base := (qMod*total + addMore[j]) % mod
                    gain := winCoeff * base % mod
                    ans[j] = (ans[j] + gain) % mod

                    nextIdx := (j + 1) % n
                    extra := addMore[(nextIdx-1+n)%n]
                    baseNext := (qMod*total + extra) % mod
                    val := winCoeff * baseNext % mod
                    nextR[currK][nextIdx] = (nextR[currK][nextIdx] + val) % mod
                }

                // Fake key case: no one wins, shift start by rem.
                if currK > 0 {
                    shift := rem % n
                    for i, v := range state {
                        if v == 0 {
                            continue
                        }
                        to := (i + shift) % n
                        next := (v * fakeProb) % mod
                        dp[currK-1][to] = (dp[currK-1][to] + next) % mod
                    }
                }
            }

            dp = nextR
        }

        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, ans[i])
        }
        fmt.Fprintln(out)
    }
}
