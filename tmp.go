package main
import (
    "fmt"
    "math/rand"
    "time"
)
func min(a, b int) int { if a < b { return a }; return b }
func solveDP(h []int) int {
    n := len(h)
    dp := make([]int, n)
    dp[0] = 0
    inc := []int{0}
    dec := []int{0}
    for i := 1; i < n; i++ {
        dp[i] = dp[i-1] + 1
        // dec stack
        for len(dec) > 0 && h[dec[len(dec)-1]] < h[i] {
            j := dec[len(dec)-1]
            dp[i] = min(dp[i], dp[j] + 1)
            dec = dec[:len(dec)-1]
        }
        if len(dec) > 0 {
            j := dec[len(dec)-1]
            dp[i] = min(dp[i], dp[j] + 1)
            if h[j] == h[i] {
                dec = dec[:len(dec)-1]
            }
        }
        dec = append(dec, i)
        // inc stack
        for len(inc) > 0 && h[inc[len(inc)-1]] > h[i] {
            j := inc[len(inc)-1]
            dp[i] = min(dp[i], dp[j] + 1)
            inc = inc[:len(inc)-1]
        }
        if len(inc) > 0 {
            j := inc[len(inc)-1]
            dp[i] = min(dp[i], dp[j] + 1)
            if h[j] == h[i] {
                inc = inc[:len(inc)-1]
            }
        }
        inc = append(inc, i)
    }
    return dp[n-1]
}
func solveBrute(h []int) int {
    n := len(h)
    INF := 1e9
    dist := make([]int, n)
    for i := 0; i < n; i++ { dist[i] = int(INF) }
    dist[0] = 0
    for i := 0; i < n; i++ {
        for j := i+1; j < n; j++ {
            ok := false
            if j == i+1 {
                ok = true
            } else {
                mx := 0
                for k := i+1; k < j; k++ { if h[k] > mx { mx = h[k] } }
                if mx < h[i] && mx < h[j] { ok = true }
                mn := 1e9
                for k := i+1; k < j; k++ { if h[k] < mn { mn = h[k] } }
                if mn > h[i] && mn > h[j] { ok = true }
            }
            if ok && dist[i] + 1 < dist[j] {
                dist[j] = dist[i] + 1
            }
        }
    }
    return dist[n-1]
}
func main() {
    rand.Seed(time.Now().UnixNano())
    for tt := 0; tt < 2000; tt++ {
        n := rand.Intn(10) + 2
        h := make([]int, n)
        for i := range h {
            h[i] = rand.Intn(5)
        }
        a := solveDP(h)
        b := solveBrute(h)
        if a != b {
            fmt.Println("Mismatch! h=", h, "dp=", a, "brute=", b)
            return
        }
    }
    fmt.Println("All tests passed")
}
