package main

import (
    "bufio"
    "fmt"
    "math/bits"
    "os"
)

const oo = 100000

// Edge represents a directed edge with capacity and cost
type Edge struct {
    to, rev, cap, cost int
}

// addEdge adds a directed edge u->v and its reverse in the residual graph
func addEdge(g [][]Edge, u, v, cap, cost int) {
    g[u] = append(g[u], Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
    g[v] = append(g[v], Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, m int
    if _, err := fmt.Fscan(in, &n, &m); err != nil {
        return
    }
    a := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(in, &a[i])
    }
    // find next equal positions
    v := make([]int, n+1)
    for i := 1; i <= n; i++ {
        j := i + 1
        for j <= n && a[j] != a[i] {
            j++
        }
        if j <= n {
            v[i] = j
        }
    }
    S := 0
    T := 3*n + 1
    tot := T + 1
    // build graph
    g := make([][]Edge, tot)
    // source to node 1 with cap m
    addEdge(g, S, 1, m, 0)
    // record edge indices for i -> i+n
    edIdx := make([]int, n+1)
    for i := 1; i <= n; i++ {
        edIdx[i] = len(g[i])
        addEdge(g, i, n+i, 1, bits.OnesCount(uint(a[i])))
        addEdge(g, n+i, 2*n+i, 1, -oo)
        addEdge(g, 2*n+i, T, 1, 0)
        if i < n {
            addEdge(g, 2*n+i, i+1, 1, 0)
            addEdge(g, i, i+1, oo, 0)
        }
    }
    // equal transitions
    for i := 1; i <= n; i++ {
        if v[i] > 0 {
            addEdge(g, 2*n+i, n+v[i], 1, 0)
        }
    }
    // min-cost flow via successive SPFA
    N := tot
    dist := make([]int, N)
    prevV := make([]int, N)
    prevE := make([]int, N)
    inQ := make([]bool, N)
    queue := make([]int, N)
    ans := 0
    for {
        for i := 0; i < N; i++ {
            dist[i] = 1<<60
            inQ[i] = false
        }
        head, tail := 0, 0
        dist[S] = 0
        queue[tail] = S
        tail++
        inQ[S] = true
        for head < tail {
            u := queue[head]
            head++
            inQ[u] = false
            for ei, e := range g[u] {
                if e.cap > 0 && dist[u]+e.cost < dist[e.to] {
                    dist[e.to] = dist[u] + e.cost
                    prevV[e.to] = u
                    prevE[e.to] = ei
                    if !inQ[e.to] {
                        inQ[e.to] = true
                        queue[tail] = e.to
                        tail++
                    }
                }
            }
        }
        if dist[T] >= 0 {
            break
        }
        ans += dist[T]
        for vtx := T; vtx != S; {
            u := prevV[vtx]
            ei := prevE[vtx]
            g[u][ei].cap--
            rev := g[u][ei].rev
            g[vtx][rev].cap++
            vtx = u
        }
    }
    // determine used edges
    used := make([]bool, n+1)
    sel := 0
    for i := 1; i <= n; i++ {
        if g[i][edIdx[i]].cap == 0 {
            used[i] = true
            sel++
        }
    }
    m2 := 2*n - sel
    cost := ans % oo
    if cost < 0 {
        cost += oo
    }
    fmt.Fprintf(out, "%d %d\n", m2, cost)
    // output operations
    now := make([]rune, n+2)
    bUsed := make([]bool, 256)
    for i := 1; i <= n; i++ {
        var w rune
        if now[i] == 0 {
            for c := 'a'; c <= 'z'; c++ {
                if !bUsed[c] {
                    w = c
                    break
                }
            }
            bUsed[w] = true
            fmt.Fprintf(out, "%c=%d\n", w, a[i])
            now[i] = w
        } else {
            w = now[i]
        }
        fmt.Fprintf(out, "print(%c)\n", w)
        if used[i] && v[i] > 0 {
            now[v[i]] = w
        } else {
            bUsed[w] = false
        }
    }
