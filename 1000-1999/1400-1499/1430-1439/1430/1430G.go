package main

import (
    "bufio"
    "fmt"
    "os"
)

const inf = 100000000

// Network implements a min-cost max-flow with successive SPFA
type Network struct {
    n    int
    to   []int
    cost []int
    cap  []int
    dist []int
    par  []int
    g    [][]int
}

// NewNetwork creates a flow network with n nodes
func NewNetwork(n int) *Network {
    g := make([][]int, n)
    return &Network{n: n, to: []int{}, cost: []int{}, cap: []int{}, dist: make([]int, n), par: make([]int, n), g: g}
}

// AddEdge adds a directed edge a->b with cost and capacity cp
func (net *Network) AddEdge(a, b, cst, cp int) {
    net.g[a] = append(net.g[a], len(net.to))
    net.to = append(net.to, b)
    net.cap = append(net.cap, cp)
    net.cost = append(net.cost, cst)
}

// Add adds edge a->b and the reverse edge b->a
func (net *Network) Add(a, b, cst, cp int) {
    net.AddEdge(a, b, cst, cp)
    net.AddEdge(b, a, -cst, 0)
}

func min(a, b int) int { if a < b { return a }; return b }

// Mincost sends 'total' flow from node 0 to node 1 minimizing cost
func (net *Network) Mincost(total int) {
    n := net.n
    inQ := make([]bool, n)
    for total > 0 {
        // initialize distances
        for i := 0; i < n; i++ {
            net.dist[i] = inf
            inQ[i] = false
        }
        net.dist[0] = 0
        // SPFA
        q := []int{0}
        inQ[0] = true
        for len(q) > 0 {
            v := q[0]
            q = q[1:]
            inQ[v] = false
            for _, it := range net.g[v] {
                if net.cap[it] > 0 {
                    u := net.to[it]
                    c := net.cost[it]
                    if net.dist[v]+c < net.dist[u] {
                        net.dist[u] = net.dist[v] + c
                        net.par[u] = it
                        if !inQ[u] {
                            q = append(q, u)
                            inQ[u] = true
                        }
                    }
                }
            }
        }
        if net.dist[1] >= inf {
            break
        }
        // find augmenting path
        mn := inf
        path := []int{}
        t := 1
        for t != 0 {
            it := net.par[t]
            path = append(path, it)
            mn = min(mn, net.cap[it])
            t = net.to[it^1]
        }
        // apply flow
        for _, it := range path {
            net.cap[it] -= mn
            net.cap[it^1] += mn
        }
        total -= mn
    }
}

func solve() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, m int
    fmt.Fscan(in, &n, &m)
    me := NewNetwork(n + 2)
    b := make([]int, n)
    for i := 0; i < m; i++ {
        var u, v, w int
        fmt.Fscan(in, &u, &v, &w)
        u--
        v--
        me.Add(u+2, v+2, -1, inf)
        b[u] += w
        b[v] -= w
    }
    total := 0
    for i := 0; i < n; i++ {
        if b[i] > 0 {
            total += b[i]
            me.Add(0, i+2, 0, b[i])
        } else if b[i] < 0 {
            me.Add(i+2, 1, 0, -b[i])
        }
    }
    me.Mincost(total)
    ans := make([]int, n)
    mn := inf
    for i := 0; i < n; i++ {
        ans[i] = me.dist[i+2]
        if ans[i] < mn {
            mn = ans[i]
        }
    }
    for i := 0; i < n; i++ {
        fmt.Fprint(out, ans[i]-mn)
        if i+1 < n {
            fmt.Fprint(out, " ")
        }
    }
}

func main() {
    solve()
}
