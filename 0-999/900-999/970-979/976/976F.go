package main

import (
   "fmt"
)

// Edge represents a flow edge
type Edge struct {
   to, rev, cap, id int
}

// Dinic holds the graph and state
type Dinic struct {
   n    int
   g    [][]Edge
   lvl  []int
   ptr  []int
   q    []int
}

// NewDinic creates a Dinic for n nodes
func NewDinic(n int) *Dinic {
   return &Dinic{
       n:   n,
       g:   make([][]Edge, n),
       lvl: make([]int, n),
       ptr: make([]int, n),
       q:   make([]int, n),
   }
}

// AddEdge adds forward and backward edges
func (d *Dinic) AddEdge(u, v, capUV, capVU, id int) {
   d.g[u] = append(d.g[u], Edge{to: v, rev: len(d.g[v]), cap: capUV, id: id})
   d.g[v] = append(d.g[v], Edge{to: u, rev: len(d.g[u]) - 1, cap: capVU, id: id})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := 0; i < d.n; i++ {
       d.lvl[i] = -1
   }
   qh, qt := 0, 0
   d.lvl[s] = 0
   d.q[qt] = s; qt++
   for qh < qt {
       v := d.q[qh]; qh++
       for _, e := range d.g[v] {
           if e.cap > 0 && d.lvl[e.to] < 0 {
               d.lvl[e.to] = d.lvl[v] + 1
               d.q[qt] = e.to; qt++
           }
       }
   }
   return d.lvl[t] >= 0
}

func (d *Dinic) dfs(v, t, f int) int {
   if v == t || f == 0 {
       return f
   }
   for i := d.ptr[v]; i < len(d.g[v]); i++ {
       e := &d.g[v][i]
       if e.cap > 0 && d.lvl[e.to] == d.lvl[v]+1 {
           pushed := d.dfs(e.to, t, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.g[e.to][e.rev].cap += pushed
               d.ptr[v] = i
               return pushed
           }
       }
   }
   d.ptr[v] = len(d.g[v])
   return 0
}

// Flow computes max flow from s to t
func (d *Dinic) Flow(s, t int) {
   for d.bfs(s, t) {
       copy(d.ptr, make([]int, d.n))
       for {
           pushed := d.dfs(s, t, 1e9)
           if pushed == 0 {
               break
           }
       }
   }
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

type EdgeEnt struct{u, v, id int}

func main() {
   var n1, n2, m int
   fmt.Scan(&n1, &n2, &m)
   N := n1 + n2
   edges := make([]EdgeEnt, m)
   currCC := make([]int, N)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Scan(&a, &b)
       a--; b--
       u := a
       v := n1 + b
       edges[i] = EdgeEnt{u: u, v: v, id: i + 1}
       currCC[u]++
       currCC[v]++
   }
   mn := currCC[0]
   for i := 1; i < N; i++ {
       if currCC[i] < mn {
           mn = currCC[i]
       }
   }
   ans := make([][]int, mn+1)
   // iterate from mn down to 0
   for i := mn; i >= 0; i-- {
       d := NewDinic(N + 2)
       S := N
       T := N + 1
       // add original edges
       origEdges := make([]*Edge, m)
       for k, e := range edges {
           d.AddEdge(e.u, e.v, 1, 0, e.id)
           // pointer to last added
           last := &d.g[e.u][len(d.g[e.u])-1]
           origEdges[k] = last
       }
       // add source/sink edges based on currCC
       for u := 0; u < N; u++ {
           if currCC[u] > i {
               if u < n1 {
                   d.AddEdge(S, u, currCC[u]-i, 0, -1)
               } else {
                   d.AddEdge(u, T, currCC[u]-i, 0, -1)
               }
           }
       }
       // max flow
       d.Flow(S, T)
       // collect kept edges and update currCC
       nextCC := make([]int, N)
       list := make([]int, 0, m)
       for k, e := range edges {
           if origEdges[k].cap > 0 {
               list = append(list, e.id)
               nextCC[e.u]++
               nextCC[e.v]++
           }
       }
       ans[i] = list
       currCC = nextCC
   }
   // output
   for i := 0; i <= mn; i++ {
       res := ans[i]
       fmt.Print(len(res))
       for _, id := range res {
           fmt.Print(" ", id)
       }
       fmt.Println()
   }
}
