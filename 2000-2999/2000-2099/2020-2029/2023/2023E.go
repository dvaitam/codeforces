package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        adj := make([][]int, n)
        edges := make([][2]int, n-1)
        for i := 0; i < n-1; i++ {
            var u, v int
            fmt.Fscan(in, &u, &v)
            u--
            v--
            adj[u] = append(adj[u], v)
            adj[v] = append(adj[v], u)
            edges[i] = [2]int{u, v}
        }

        deg := make([]int, n)
        for i := 0; i < n; i++ {
            deg[i] = len(adj[i])
        }

        var totalPairs int64
        for i := 0; i < n; i++ {
            d := int64(deg[i])
            totalPairs += d * (d - 1) / 2
        }

        var pairMerges int64
        for _, e := range edges {
            du := int64(deg[e[0]] - 1)
            dv := int64(deg[e[1]] - 1)
            if du < 0 {
                du = 0
            }
            if dv < 0 {
                dv = 0
            }
            if du < dv {
                pairMerges += du
            } else {
                pairMerges += dv
            }
        }

        extra := int64(0)
        seen := make(map[uint64]struct{})
        for v := 0; v < n; v++ {
            if deg[v] < 3 {
                continue
            }
            for _, to := range adj[v] {
                end := to
                prev := v
                length := 1
                for deg[end] == 2 {
                    next := adj[end][0]
                    if next == prev {
                        next = adj[end][1]
                    }
                    prev = end
                    end = next
                    length++
                }
                if deg[end] >= 3 && length >= 2 {
                    a, b := v, end
                    if a > b {
                        a, b = b, a
                    }
                    key := (uint64(a) << 32) | uint64(b)
                    if _, ok := seen[key]; !ok {
                        seen[key] = struct{}{}
                        extra++
                    }
                }
            }
        }

        ans := totalPairs - pairMerges - extra
        fmt.Fprintln(out, ans)
    }
}
