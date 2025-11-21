package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const MOD int64 = 998244353

var fact []int64 = []int64{1}
var invFact []int64 = []int64{1}

func modPow(a int64, e int) int64 {
    if a %= MOD; a < 0 {
        a += MOD
    }
    res := int64(1)
    base := a
    exp := e
    for exp > 0 {
        if exp&1 == 1 {
            res = res * base % MOD
        }
        base = base * base % MOD
        exp >>= 1
    }
    return res
}

func modInverse(a int64) int64 {
    return modPow(a, int(MOD-2))
}

func ensureFact(n int) {
    if n < 0 {
        return
    }
    if len(fact) > n {
        return
    }
    old := len(fact) - 1
    need := n - old
    fact = append(fact, make([]int64, need)...)
    for i := old + 1; i <= n; i++ {
        fact[i] = fact[i-1] * int64(i) % MOD
    }
    invFact = append(invFact, make([]int64, need)...)
    invFact[n] = modInverse(fact[n])
    for i := n; i > old+0; i-- {
        invFact[i-1] = invFact[i] * int64(i) % MOD
    }
}

func segmentWays(l, r, B int) int64 {
    if l > r {
        return 1
    }
    if B < 0 {
        return 0
    }
    split := B + 1
    first := r
    if first > split {
        first = split
    }
    prod1 := int64(1)
    if first >= l {
        ensureFact(first)
        prod1 = fact[first] * invFact[l-1] % MOD
    }
    last := l
    if last < split+1 {
        last = split + 1
    }
    count := 0
    if last <= r {
        count = r - last + 1
    }
    prod2 := modPow(int64(B+1), count)
    return prod1 * prod2 % MOD
}

func readInt(r *bufio.Reader) int {
    sign := 1
    val := 0
    c, _ := r.ReadByte()
    for (c < '0' || c > '9') && c != '-' {
        c, _ = r.ReadByte()
    }
    if c == '-' {
        sign = -1
        c, _ = r.ReadByte()
    }
    for c >= '0' && c <= '9' {
        val = val*10 + int(c-'0')
        c, _ = r.ReadByte()
    }
    return sign * val
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    t := readInt(in)
    for ; t > 0; t-- {
        n := readInt(in)
        m := readInt(in)

        bounds := make(map[int]struct{})
        bounds[0] = struct{}{}
        bounds[n] = struct{}{}
        upper := make(map[int]int)
        lower := make(map[int]int)
        valsSet := make(map[int]struct{})
        valsSet[-1] = struct{}{}
        valsSet[0] = struct{}{}
        valsSet[n-1] = struct{}{}

        type Query struct{ k, l, r int }
        queries := make([]Query, m)
        for i := 0; i < m; i++ {
            k := readInt(in)
            l := readInt(in)
            r := readInt(in)
            queries[i] = Query{k, l, r}
            bounds[l] = struct{}{}
            if cur, ok := upper[l]; ok {
                if k < cur {
                    upper[l] = k
                }
            } else {
                upper[l] = k
            }
            vv := k
            if vv > n-1 {
                vv = n - 1
            }
            if vv < 0 {
                vv = 0
            }
            valsSet[vv] = struct{}{}
            vv2 := k + 1
            if vv2 > n-1 {
                vv2 = n - 1
            }
            if vv2 < 0 {
                vv2 = 0
            }
            valsSet[vv2] = struct{}{}
            if r < n {
                pos := r + 1
                bounds[pos] = struct{}{}
                need := k + 1
                if cur, ok := lower[pos]; ok {
                    if need > cur {
                        lower[pos] = need
                    }
                } else {
                    lower[pos] = need
                }
                vv3 := k + 1
                if vv3 > n-1 {
                    vv3 = n - 1
                }
                if vv3 < 0 {
                    vv3 = 0
                }
                valsSet[vv3] = struct{}{}
            }
        }

        posList := make([]int, 0, len(bounds))
        for p := range bounds {
            posList = append(posList, p)
        }
        sort.Ints(posList)

        L := make([]int, len(posList))
        U := make([]int, len(posList))
        impossible := false
        for i, p := range posList {
            cap := 0
            if p > 0 {
                cap = p - 1
            }
            if cap > n-1 {
                cap = n - 1
            }
            U[i] = cap
            if val, ok := upper[p]; ok {
                if val < U[i] {
                    U[i] = val
                }
            }
            if val, ok := lower[p]; ok {
                if val > L[i] {
                    L[i] = val
                }
            }
            if L[i] < 0 {
                L[i] = 0
            }
            if U[i] > n-1 {
                U[i] = n - 1
            }
            if L[i] > n-1 {
                impossible = true
                break
            }
            if U[i] < 0 {
                U[i] = 0
            }
            if L[i] > U[i] {
                impossible = true
                break
            }
            valsSet[U[i]] = struct{}{}
            valsSet[L[i]] = struct{}{}
        }
        if impossible {
            fmt.Fprintln(out, 0)
            continue
        }

        vals := make([]int, 0, len(valsSet))
        for v := range valsSet {
            if v <= n-1 {
                vals = append(vals, v)
            }
        }
        sort.Ints(vals)
        if len(vals) == 0 || vals[0] != -1 {
            vals = append([]int{-1}, vals...)
        }

        ensureFact(n)

        dp := make([]int64, len(vals))
        j0 := 1
        for j := 1; j < len(vals); j++ {
            if vals[j] <= 0 {
                j0 = j
            } else {
                break
            }
        }
        dp[j0] = 1

        for idx := 1; idx < len(posList); idx++ {
            l := posList[idx-1] + 1
            r := posList[idx]
            g := make([]int64, len(vals))
            for j := 0; j < len(vals); j++ {
                g[j] = segmentWays(l, r, vals[j])
            }
            newDP := make([]int64, len(vals))
            pref := int64(0)
            for j := 1; j < len(vals); j++ {
                stay := dp[j] * g[j] % MOD
                delta := g[j] - g[j-1]
                if delta < 0 {
                    delta += MOD
                }
                rise := pref * delta % MOD
                newDP[j] = (stay + rise) % MOD
                pref = (pref + dp[j]) % MOD
            }
            low := L[idx]
            high := U[idx]
            for j := 1; j < len(vals); j++ {
                if vals[j] < low || vals[j] > high {
                    newDP[j] = 0
                }
            }
            dp = newDP
        }

        ans := int64(0)
        for j := 1; j < len(vals); j++ {
            ans += dp[j]
            if ans >= MOD {
                ans -= MOD
            }
        }
        fmt.Fprintln(out, ans%MOD)
    }
}
