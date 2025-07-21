package main

import (
    "fmt"
)

func main() {
    var s1, s2 string
    var n int
    if _, err := fmt.Scan(&s1, &s2, &n); err != nil {
        return
    }
    // Read productions a->bc
    type prod struct{ a, b, c int }
    prods := make([]prod, 0, n)
    for i := 0; i < n; i++ {
        var rule string
        fmt.Scan(&rule)
        if len(rule) == 5 && rule[1] == '-' && rule[2] == '>' {
            A := int(rule[0] - 'a')
            B := int(rule[3] - 'a')
            C := int(rule[4] - 'a')
            prods = append(prods, prod{A, B, C})
        }
    }
    // DP masks: dp1[i][j] bitmask of letters that can derive s1[i:j] inclusive
    n1 := len(s1)
    n2 := len(s2)
    dp1 := make([][]uint32, n1)
    for i := range dp1 {
        dp1[i] = make([]uint32, n1)
    }
    dp2 := make([][]uint32, n2)
    for i := range dp2 {
        dp2[i] = make([]uint32, n2)
    }
    // base case
    for i := 0; i < n1; i++ {
        dp1[i][i] = 1 << (s1[i] - 'a')
    }
    for i := 0; i < n2; i++ {
        dp2[i][i] = 1 << (s2[i] - 'a')
    }
    // build DP for s1
    for length := 2; length <= n1; length++ {
        for i := 0; i+length <= n1; i++ {
            j := i + length - 1
            var mask uint32
            for k := i; k < j; k++ {
                left := dp1[i][k]
                right := dp1[k+1][j]
                if left == 0 || right == 0 {
                    continue
                }
                for _, p := range prods {
                    if (left>>p.b)&1 != 0 && (right>>p.c)&1 != 0 {
                        mask |= 1 << p.a
                    }
                }
            }
            dp1[i][j] = mask
        }
    }
    // build DP for s2
    for length := 2; length <= n2; length++ {
        for i := 0; i+length <= n2; i++ {
            j := i + length - 1
            var mask uint32
            for k := i; k < j; k++ {
                left := dp2[i][k]
                right := dp2[k+1][j]
                if left == 0 || right == 0 {
                    continue
                }
                for _, p := range prods {
                    if (left>>p.b)&1 != 0 && (right>>p.c)&1 != 0 {
                        mask |= 1 << p.a
                    }
                }
            }
            dp2[i][j] = mask
        }
    }
    // BFS on grid of positions
    const INF = 1e9
    dist := make([][]int, n1+1)
    for i := 0; i <= n1; i++ {
        dist[i] = make([]int, n2+1)
        for j := range dist[i] {
            dist[i][j] = INF
        }
    }
    type pair struct{ i, j int }
    q := make([]pair, 0, (n1+1)*(n2+1))
    dist[0][0] = 0
    q = append(q, pair{0, 0})
    for head := 0; head < len(q); head++ {
        u := q[head]
        d := dist[u.i][u.j]
        // try all segment ends
        for i2 := u.i + 1; i2 <= n1; i2++ {
            for j2 := u.j + 1; j2 <= n2; j2++ {
                m1 := dp1[u.i][i2-1]
                if m1 == 0 {
                    continue
                }
                m2 := dp2[u.j][j2-1]
                if m2 == 0 {
                    continue
                }
                if (m1 & m2) == 0 {
                    continue
                }
                if dist[i2][j2] > d+1 {
                    dist[i2][j2] = d + 1
                    q = append(q, pair{i2, j2})
                }
            }
        }
    }
    ans := dist[n1][n2]
    if ans >= INF {
        fmt.Println(-1)
    } else {
        fmt.Println(ans)
    }
}
