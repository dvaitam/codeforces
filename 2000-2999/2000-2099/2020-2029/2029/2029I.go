package main

import (
    "bufio"
    "fmt"
    "os"
)

const inf int64 = 1<<60

// We model the number of operations as the number of interval starts.
// For a coverage array c, the minimal number of intervals needed equals
// c[0] + sum max(0, c[i]-c[i-1]). Splitting an interval does not change
// coverage, so any count p in [minimal, totalLength] is feasible.
//
// DP state: for a fixed number of starts `u`, we keep a map keyed by
// (totalCoverageSum, currentHeight) and store the minimal additive part of
// the objective: k^2*n*sum c_i^2 + 2k*sum (n*a_i-A)*c_i.
//
// The global term -k^2*(sum c_i)^2 is recovered when evaluating answers.
// Because n*m <= 2e4, the number of reachable (sum,height) pairs is
// manageable with sparse maps even in the worst case (small n or small m).
//
// Key encoding packs sum and height to a single int64 to index the maps.
func main() {
    in := bufio.NewReader(os.Stdin)
    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    for ; T > 0; T-- {
        var n, m int
        var k int64
        fmt.Fscan(in, &n, &m, &k)
        a := make([]int64, n)
        var sumA, sumA2 int64
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            sumA += a[i]
            sumA2 += a[i] * a[i]
        }

        // Precompute constants.
        base := int64(n)*sumA2 - sumA*sumA
        w := make([]int64, n)
        for i := 0; i < n; i++ {
            w[i] = int64(n)*a[i] - sumA
        }

        maxLen := n * m // upper bound on total coverage length
        keyMod := int64(m + 1)
        encode := func(sum int64, h int) int64 {
            return sum*keyMod + int64(h)
        }

        // cur[u] holds map[key] -> cost for states after processing i positions.
        cur := make([]map[int64]int64, m+1)
        cur[0] = map[int64]int64{encode(0, 0): 0}

        k2 := k * k
        k2n := k2 * int64(n)

        for idx := 0; idx < n; idx++ {
            // Precompute per-height cost for this position.
            addCost := make([]int64, m+1)
            for h := 0; h <= m; h++ {
                addCost[h] = k2n*int64(h*h) + 2*k*w[idx]*int64(h)
            }

            next := make([]map[int64]int64, m+1)

            for u := 0; u <= m; u++ {
                if cur[u] == nil {
                    continue
                }
                for key, val := range cur[u] {
                    sum := key / keyMod
                    h := int(key % keyMod)
                    for h2 := 0; h2 <= m; h2++ {
                        inc := h2 - h
                        if inc < 0 {
                            inc = 0
                        }
                        u2 := u + inc
                        if u2 > m {
                            continue
                        }
                        sum2 := sum + int64(h2)
                        if sum2 > int64(maxLen) {
                            continue
                        }
                        cost2 := val + addCost[h2]
                        if next[u2] == nil {
                            next[u2] = make(map[int64]int64)
                        }
                        k2key := encode(sum2, h2)
                        if old, ok := next[u2][k2key]; !ok || cost2 < old {
                            next[u2][k2key] = cost2
                        }
                    }
                }
            }
            cur = next
        }

        ans := make([]int64, m+1) // 1-based for convenience
        for i := 1; i <= m; i++ {
            ans[i] = inf
        }

        for u := 0; u <= m; u++ {
            mp := cur[u]
            if mp == nil {
                continue
            }
            for key, cost := range mp {
                sum := key / keyMod
                if sum == 0 {
                    continue
                }
                val := base + cost - k2*sum*sum
                maxP := m
                if sum < int64(maxP) {
                    maxP = int(sum)
                }
                if u > maxP {
                    continue
                }
                for p := u; p <= maxP; p++ {
                    if val < ans[p] {
                        ans[p] = val
                    }
                }
            }
        }

        // Fill remaining answers with previous (non-increasing property)
        for p := 1; p <= m; p++ {
            if p > 1 && ans[p] > ans[p-1] {
                ans[p] = ans[p-1]
            }
        }

        for p := 1; p <= m; p++ {
            if p > 1 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, ans[p])
        }
        fmt.Fprintln(out)
    }
}
