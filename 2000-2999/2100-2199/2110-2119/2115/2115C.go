package main

import (
    "bufio"
    "fmt"
    "os"
)

type state struct {
    r int // rounds remaining
    d int // total single-damage still needed
    s int // shines that can still be safely used
}

var (
    memo map[state]float64
    probShine float64
    n int
)

// Maximum damage we can still deal in r rounds with s shine opportunities.
// Each shine round contributes n damage instead of 1, and we can shine at
// most once per round.
func maxDamage(r, s int) int {
    if s > r {
        s = r
    }
    return r + (n-1)*s
}

func solve(r, d, s int) float64 {
    if d == 0 {
        return 1.0
    }
    if r == 0 {
        return 0.0
    }
    if d > maxDamage(r, s) {
        return 0.0
    }

    key := state{r, d, s}
    if v, ok := memo[key]; ok {
        return v
    }

    // Sword shines this round.
    var nextShine float64
    if s > 0 && d >= n {
        nextShine = solve(r-1, d-n, s-1)
    } else {
        // Either no shine left or using it would overkill, so skip.
        nextShine = solve(r-1, d, s)
    }

    // Sword does not shine: choose optimally between single attack and skip.
    bestNonShine := 0.0
    // Option 1: skip this round (only if still feasible afterwards).
    if d <= maxDamage(r-1, s) {
        cand := solve(r-1, d, s)
        if cand > bestNonShine {
            bestNonShine = cand
        }
    }
    // Option 2: perform a single-target attack if there is damage left.
    if d > 0 {
        cand := solve(r-1, d-1, s)
        if cand > bestNonShine {
            bestNonShine = cand
        }
    }

    res := probShine*nextShine + (1-probShine)*bestNonShine
    memo[key] = res
    return res
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var m int
        var pInt int
        fmt.Fscan(in, &n, &m, &pInt)
        h := make([]int, n)
        minH := 1 << 30
        total := 0
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &h[i])
            need := h[i] - 1
            total += need
            if need < minH {
                minH = need
            }
        }

        probShine = float64(pInt) / 100.0
        memo = make(map[state]float64)

        ans := solve(m, total, minH)
        fmt.Fprintf(out, "%.9f\n", ans)
    }
}
