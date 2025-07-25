package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000

// Edge for Dinic
type Edge struct {
   to, rev, cap int
}

// Dinic structure
type Dinic struct {
   n     int
   graph [][]Edge
   level []int
   ptr   []int
   q     []int
}

func NewDinic(n int) *Dinic {
   return &Dinic{
       n:     n,
       graph: make([][]Edge, n),
       level: make([]int, n),
       ptr:   make([]int, n),
       q:     make([]int, n),
   }
}

func (d *Dinic) AddEdge(u, v, c int) {
   d.graph[u] = append(d.graph[u], Edge{to: v, rev: len(d.graph[v]), cap: c})
   d.graph[v] = append(d.graph[v], Edge{to: u, rev: len(d.graph[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := 0; i < d.n; i++ {
       d.level[i] = -1
   }
   ql, qr := 0, 0
   d.q[qr] = s; qr++
   d.level[s] = 0
   for ql < qr {
       u := d.q[ql]; ql++
       for _, e := range d.graph[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               d.q[qr] = e.to; qr++
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
   if u == t {
       return f
   }
   for i := d.ptr[u]; i < len(d.graph[u]); i++ {
       e := &d.graph[u][i]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           pushed := d.dfs(e.to, t, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.graph[e.to][e.rev].cap += pushed
               return pushed
           }
       }
       d.ptr[u]++
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := 0; i < d.n; i++ {
           d.ptr[i] = 0
       }
       for {
           pushed := d.dfs(s, t, INF)
           if pushed == 0 {
               break
           }
           flow += pushed
       }
   }
   return flow
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, h, m int
   if _, err := fmt.Fscan(reader, &n, &h, &m); err != nil {
       return
   }
   // IDs: 0 = source, 1 = sink
   // b-nodes: for i=1..n, k=1..h: id = 2 + (i-1)*h + (k-1)
   baseB := 2
   baseR := baseB + n*h
   totalNodes := baseR + m
   d := NewDinic(totalNodes)
   src, sink := 0, 1
   // add b-node edges
   for i := 1; i <= n; i++ {
       for k := 1; k <= h; k++ {
           id := baseB + (i-1)*h + (k-1)
           w := 2*k - 1
           d.AddEdge(src, id, w)
           if k > 1 {
               prev := baseB + (i-1)*h + (k-2)
               d.AddEdge(id, prev, INF)
           }
       }
   }
   // restrictions
   for j := 0; j < m; j++ {
       var l, r, x, c int
       fmt.Fscan(reader, &l, &r, &x, &c)
       rnode := baseR + j
       if x < h {
           k := x + 1
           for i := l; i <= r; i++ {
               bid := baseB + (i-1)*h + (k-1)
               d.AddEdge(bid, rnode, INF)
           }
           d.AddEdge(rnode, sink, c)
       }
       // if x >= h, ignore
   }
   // total potential profit = n * h^2
   total := n * h * h
   flow := d.MaxFlow(src, sink)
   fmt.Println(total - flow)
}
