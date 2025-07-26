package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod = 1000000007
const maxN = 200
const maxC = maxN*(maxN+1)/2

var comb [maxN+1][maxN+1]int
var dp [][]int

func initComb() {
    for i := 0; i <= maxN; i++ {
        comb[i][0] = 1
        comb[i][i] = 1
        for j := 1; j < i; j++ {
            comb[i][j] = (comb[i-1][j-1] + comb[i-1][j]) % mod
        }
    }
}

func compute(n int) {
    dp = make([][]int, n+1)
    dp[0] = []int{1}
    for i := 1; i <= n; i++ {
        dp[i] = make([]int, i*(i+1)/2+1)
        for l := 0; l < i; l++ {
            r := i - 1 - l
            coef := comb[i-1][l]
            for a, va := range dp[l] {
                if va == 0 {
                    continue
                }
                for b, vb := range dp[r] {
                    if vb == 0 {
                        continue
                    }
                    c := i + a + b
                    val := (int64(va) * int64(vb)) % mod
                    val = (val * int64(coef)) % mod
                    dp[i][c] = int((int64(dp[i][c]) + val) % mod)
                }
            }
        }
    }
}

func main() {
    in := bufio.NewReader(os.Stdin)
    var n, lim int
    fmt.Fscan(in, &n, &lim)
    maxCost := n * (n + 1) / 2
    if lim > maxCost {
        fmt.Println(0)
        return
    }
    initComb()
    compute(n)
    ans := 0
    for c := lim; c <= maxCost; c++ {
        ans += dp[n][c]
        if ans >= mod {
            ans -= mod
        }
    }
    fmt.Println(ans)
}

