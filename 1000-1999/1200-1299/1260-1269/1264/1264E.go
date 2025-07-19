package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

// edge represents a flow edge
type edge struct {
   to, rev, cap, cost, orig int
}

// MCMF holds graph for min cost max flow
type MCMF struct {
   n       int
   graph   [][]edge
   dist    []int
   prevv   []int
   preve   []int
}

// NewMCMF creates MCMF with n nodes
func NewMCMF(n int) *MCMF {
   g := make([][]edge, n)
   return &MCMF{n: n, graph: g, dist: make([]int, n), prevv: make([]int, n), preve: make([]int, n)}
}

// AddEdge adds directed edge u->v
func (f *MCMF) AddEdge(u, v, cap, cost int) {
   f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost, orig: cap})
   f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost, orig: 0})
}

// MinCostFlow runs successive SPFA to compute min cost max flow
func (f *MCMF) MinCostFlow(s, t int) int {
   res := 0
   for {
       // SPFA
       for i := 0; i < f.n; i++ {
           f.dist[i] = INF
       }
       inq := make([]bool, f.n)
       queue := make([]int, 0, f.n)
       f.dist[s] = 0
       inq[s] = true
       queue = append(queue, s)
       for idx := 0; idx < len(queue); idx++ {
           u := queue[idx]
           inq[u] = false
           for i, e := range f.graph[u] {
               if e.cap > 0 && f.dist[e.to] > f.dist[u] + e.cost {
                   f.dist[e.to] = f.dist[u] + e.cost
                   f.prevv[e.to] = u
                   f.preve[e.to] = i
                   if !inq[e.to] {
                       queue = append(queue, e.to)
                       inq[e.to] = true
                   }
               }
           }
       }
       if f.dist[t] == INF {
           break
       }
       // add as much as possible (here caps are 1)
       d := INF
       for v := t; v != s; v = f.prevv[v] {
           e := f.graph[f.prevv[v]][f.preve[v]]
           if d > e.cap {
               d = e.cap
           }
       }
       for v := t; v != s; v = f.prevv[v] {
           e := &f.graph[f.prevv[v]][f.preve[v]]
           e.cap -= d
           f.graph[v][e.rev].cap += d
       }
       res += d * f.dist[t]
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n)
   // rel holds initial relations
   rel := make([][]int, n+1)
   for i := range rel {
       rel[i] = make([]int, n+1)
       for j := range rel[i] {
           rel[i][j] = 2
       }
   }
   indeg := make([]int, n+1)
   fmt.Fscan(in, &m)
   for k := 0; k < m; k++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       rel[x][y] = 1
       rel[y][x] = 0
       indeg[y]++
   }
   // assign pair node ids
   pairID := make([][]int, n+1)
   for i := range pairID {
       pairID[i] = make([]int, n+1)
   }
   id := n
   for i := 1; i <= n; i++ {
       for j := i + 1; j <= n; j++ {
           id++
           pairID[i][j] = id
       }
   }
   t := id + 1
   s := 0
   mcmf := NewMCMF(t + 1)
   // wedge map for flows
   type wi struct{from, idx int}
   wedge := make([][]wi, n+1)
   for i := range wedge {
       wedge[i] = make([]wi, n+1)
   }
   // add pair edges
   for i := 1; i <= n; i++ {
       for j := i + 1; j <= n; j++ {
           if rel[i][j] < 2 {
               continue
           }
           pid := pairID[i][j]
           mcmf.AddEdge(s, pid, 1, 0)
           // to i
           wk := len(mcmf.graph[pid])
           mcmf.AddEdge(pid, i, 1, 0)
           wedge[j][i] = wi{from: pid, idx: wk}
           // to j
           wk2 := len(mcmf.graph[pid])
           mcmf.AddEdge(pid, j, 1, 0)
           wedge[i][j] = wi{from: pid, idx: wk2}
       }
   }
   // edges from nodes to t
   for i := 1; i <= n; i++ {
       for j := indeg[i] + 1; j < n; j++ {
           mcmf.AddEdge(i, t, 1, j-1)
       }
   }
   // run flow
   _ = mcmf.MinCostFlow(s, t)
   // print result
   relOut := rel // rel matrix updated
   for i := 1; i <= n; i++ {
       relOut[i][i] = 0
       for j := 1; j <= n; j++ {
           if relOut[i][j] < 2 {
               fmt.Print(relOut[i][j])
           } else {
               w := wedge[i][j]
               e := mcmf.graph[w.from][w.idx]
               flow := e.orig - e.cap
               fmt.Print(flow)
           }
       }
       fmt.Println()
   }
}
