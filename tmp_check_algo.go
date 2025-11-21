package main
import (
    "fmt"
    "math/rand"
    "time"
)

// brute
func brute(a []int64, b []int) int64 {
    memo := map[int]map[int]int64{}
    var solve func(mask int, cur int) int64
    solve = func(mask int, cur int) int64 {
        if memo[mask] == nil {
            memo[mask] = map[int]int64{}
        }
        if v, ok := memo[mask][cur]; ok {
            return v
        }
        idx := cur - 1
        best := int64(0)
        // submit
        nmask := mask | (1 << idx)
        next := nextIdx(nmask, cur-1, len(a))
        cand := a[idx]
        if next != 0 {
            cand += solve(nmask, next)
        }
        if cand > best {
            best = cand
        }
        // skip
        nb := b[idx]
        next = nextIdx(nmask, nb, len(a))
        cand = 0
        if next != 0 {
            cand += solve(nmask, next)
        }
        if cand > best {
            best = cand
        }

        memo[mask][cur] = best
        return best
    }
    return solve(0, 1)
}

func nextIdx(mask int, bound int, n int) int {
    for i := bound; i >= 1; i-- {
        if mask&(1<<(i-1)) == 0 {
            return i
        }
    }
    return 0
}

const INF int64 = 1<<60

func fast(a []int64, b []int) int64 {
    n := len(a)
    // compute R
    R := 1
    for i := 1; i <= R; i++ {
        if b[i-1] > R {
            R = b[i-1]
        }
    }
    if R > n {
        R = n
    }
    pref := make([]int64, n+1)
    for i, v := range a {
        pref[i+1] = pref[i] + v
    }
    if R == 1 {
        if a[0] > 0 {
            return a[0]
        }
        return 0
    }
    size := R - 1
    seg := make([]int64, 4*size+5)
    for i := range seg {
        seg[i] = INF
    }
    var upd func(node, l, r, ql, qr int, val int64)
    upd = func(node, l, r, ql, qr int, val int64) {
        if ql <= l && r <= qr {
            if val < seg[node] {
                seg[node] = val
            }
            return
        }
        m := (l + r) >> 1
        if ql <= m {
            upd(node<<1, l, m, ql, qr, val)
        }
        if qr > m {
            upd(node<<1|1, m+1, r, ql, qr, val)
        }
    }
    var query func(node, l, r, idx int) int64
    query = func(node, l, r, idx int) int64 {
        res := seg[node]
        if l == r {
            return res
        }
        m := (l + r) >> 1
        if idx <= m {
            t := query(node<<1, l, m, idx)
            if t < res {
                res = t
            }
            return res
        }
        t := query(node<<1|1, m+1, r, idx)
        if t < res {
            res = t
        }
        return res
    }

    dp := make([]int64, R)
    for i := range dp {
        dp[i] = INF
    }
    dp[0] = 0

    for i := 1; i <= R; i++ {
        val := dp[i-1] + a[i-1]
        l := i
        r := b[i-1] - 1
        if r > R-1 {
            r = R - 1
        }
        if l <= r {
            upd(1, 1, size, l, r, val)
        }
        if i <= R-1 {
            dp[i] = query(1, 1, size, i)
        }
    }

    ans := int64(0)
    for k := 1; k <= R; k++ {
        if dp[k-1] >= INF {
            continue
        }
        cur := pref[k] - dp[k-1]
        if cur > ans {
            ans = cur
        }
    }
    return ans
}

func main() {
    rand.Seed(time.Now().UnixNano())
    for n := 1; n <= 10; n++ {
        for t := 0; t < 10000; t++ {
            a := make([]int64, n)
            b := make([]int, n)
            for i := 0; i < n; i++ {
                a[i] = int64(rand.Intn(5) + 1)
                b[i] = rand.Intn(n) + 1
            }
            br := brute(a, b)
            fs := fast(a, b)
            if br != fs {
                fmt.Println("Mismatch n", n, "a", a, "b", b, "br", br, "fast", fs)
                return
            }
        }
        fmt.Println("n", n, "ok")
    }
    fmt.Println("all ok")
}
