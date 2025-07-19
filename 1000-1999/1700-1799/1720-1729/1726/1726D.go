package main

import (
    "bufio"
    "fmt"
    "os"
)

// DSU with path compression
type DSU struct {
    p []int
}

func NewDSU(n int) *DSU {
    p := make([]int, n+1)
    for i := 1; i <= n; i++ {
        p[i] = i
    }
    return &DSU{p: p}
}

func (d *DSU) find(x int) int {
    if d.p[x] != x {
        d.p[x] = d.find(d.p[x])
    }
    return d.p[x]
}

// union returns true if merged, false if already in same set
func (d *DSU) union(a, b int) bool {
    fa := d.find(a)
    fb := d.find(b)
    if fa == fb {
        return false
    }
    d.p[fa] = fb
    return true
}

// pair represents an adjacency entry: neighbor and edge index
type pair struct {
    to, idx int
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for T > 0 {
        T--
        var n, m int
        fmt.Fscan(in, &n, &m)
        // answer flags for edges
        ans := make([]byte, m)
        for i := 0; i < m; i++ {
            ans[i] = '0'
        }
        edges := make([][2]int, m)
        for i := 0; i < m; i++ {
            fmt.Fscan(in, &edges[i][0], &edges[i][1])
        }

        dsu := NewDSU(n)
        adj := make([][]pair, n+1)
        extra := make([]int, 0, 3)

        for i := 0; i < m; i++ {
            a := edges[i][0]
            b := edges[i][1]
            if dsu.union(a, b) {
                // include in spanning forest
                ans[i] = '1'
                adj[a] = append(adj[a], pair{to: b, idx: i})
                adj[b] = append(adj[b], pair{to: a, idx: i})
            } else {
                extra = append(extra, i)
            }
        }

        // special case: one cycle of length 3
        if m == n+2 && len(extra) == 3 {
            nodes := make(map[int]struct{}, 3)
            for _, i := range extra {
                nodes[edges[i][0]] = struct{}{}
                nodes[edges[i][1]] = struct{}{}
            }
            if len(nodes) == 3 {
                // pick first extra edge to include
                d0 := extra[0]
                ans[d0] = '1'
                // remove adjacent spanning edges at one endpoint
                p := edges[d0][0]
                for _, pr := range adj[p] {
                    ans[pr.idx] = '0'
                }
            }
        }

        out.Write(ans)
        out.WriteByte('\n')
    }
}
