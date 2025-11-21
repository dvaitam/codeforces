package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod int64 = 998244353

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
    in := bufio.NewReader(os.Stdin)
    var n, m int
    if _, err := fmt.Fscan(in, &n, &m); err != nil {
        return
    }
    a := make([]int64, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
    }
    b := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &b[i])
    }

    limit := n + 5
    fac := make([]int64, limit)
    ifac := make([]int64, limit)
    fac[0] = 1
    for i := 1; i < limit; i++ {
        fac[i] = fac[i-1] * int64(i) % mod
    }
    ifac[limit-1] = modPow(fac[limit-1], mod-2)
    for i := limit - 2; i >= 0; i-- {
        ifac[i] = ifac[i+1] * int64(i+1) % mod
    }

    probs := make([][]int64, n)
    for i := range probs {
        probs[i] = make([]int64, m+1)
    }

    for i := n - 1; i >= 0; i-- {
        if b[i] > 0 {
            probs[i][b[i]] = 1
            continue
        }
        type entry struct {
            idx int
            val int64
        }
        var data []entry
        firstInit := -1
        var firstInitVal int64
        for j := i + 1; j < n; j++ {
            if a[j]%a[i] != 0 {
                continue
            }
            if b[j] > 0 {
                if firstInit == -1 || a[j] < firstInitVal {
                    firstInit = j
                    firstInitVal = a[j]
                }
            } else {
                data = append(data, entry{idx: j, val: a[j]})
            }
        }
        var cand []int
        if firstInit != -1 {
            for _, e := range data {
                if e.val < firstInitVal {
                    cand = append(cand, e.idx)
                } else {
                    break
                }
            }
        } else {
            cand = make([]int, len(data))
            for k, e := range data {
                cand[k] = e.idx
            }
        }

        u := len(cand)
        eArr := make([]int64, u+1)
        eArr[0] = 1
        for pos, idx := range cand {
            invDen := ifac[pos+2]
            var weight int64
            for t := 0; t <= pos; t++ {
                coeff := fac[t+1]
                coeff = coeff * fac[pos-t] % mod
                coeff = coeff * invDen % mod
                weight = (weight + coeff*eArr[t]) % mod
            }
            row := probs[idx]
            for s := 1; s <= m; s++ {
                if row[s] == 0 {
                    continue
                }
                probs[i][s] = (probs[i][s] + weight*row[s]) % mod
            }
            discard := row[0]
            if discard != 0 {
                for t := pos; t >= 0; t-- {
                    eArr[t+1] = (eArr[t+1] + eArr[t]*discard) % mod
                }
            }
        }
        invDen := ifac[u+1]
        var tail int64
        for t := 0; t <= u; t++ {
            coeff := fac[t]
            coeff = coeff * fac[u-t] % mod
            coeff = coeff * invDen % mod
            tail = (tail + coeff*eArr[t]) % mod
        }
        if firstInit != -1 {
            owner := b[firstInit]
            probs[i][owner] = (probs[i][owner] + tail) % mod
        } else {
            probs[i][0] = (probs[i][0] + tail) % mod
        }
    }

    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    for s := 1; s <= m; s++ {
        var ans int64
        for i := 0; i < n; i++ {
            ans = (ans + probs[i][s]) % mod
        }
        if s > 1 {
            fmt.Fprint(out, " ")
        }
        fmt.Fprint(out, ans)
    }
    fmt.Fprintln(out)
}
