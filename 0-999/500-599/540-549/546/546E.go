package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a directed edge with capacity and a reverse index
type Edge struct {
   to, cap, rev int
}

// Dinic implements Dinic's max flow algorithm
type Dinic struct {
   N     int
   G     [][]Edge
   level []int
   iter  []int
}

// NewDinic creates a new Dinic struct with n nodes
func NewDinic(n int) *Dinic {
   G := make([][]Edge, n)
   return &Dinic{N: n, G: G, level: make([]int, n), iter: make([]int, n)}
}

// AddEdge adds a directed edge u->v with capacity c
func (d *Dinic) AddEdge(u, v, c int) {
   d.G[u] = append(d.G[u], Edge{to: v, cap: c, rev: len(d.G[v])})
   d.G[v] = append(d.G[v], Edge{to: u, cap: 0, rev: len(d.G[u]) - 1})
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// bfs builds level graph
func (d *Dinic) bfs(s int) {
   for i := range d.level {
       d.level[i] = -1
   }
   queue := make([]int, 0, d.N)
   d.level[s] = 0
   queue = append(queue, s)
   for i := 0; i < len(queue); i++ {
       v := queue[i]
       for _, e := range d.G[v] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[v] + 1
               queue = append(queue, e.to)
           }
       }
   }
}

// dfs finds an augmenting path using DFS on level graph
func (d *Dinic) dfs(v, t, f int) int {
   if v == t {
       return f
   }
   for ; d.iter[v] < len(d.G[v]); d.iter[v]++ {
       e := &d.G[v][d.iter[v]]
       if e.cap > 0 && d.level[v] < d.level[e.to] {
           ret := d.dfs(e.to, t, min(f, e.cap))
           if ret > 0 {
               e.cap -= ret
               d.G[e.to][e.rev].cap += ret
               return ret
           }
       }
   }
   return 0
}

// MaxFlow computes the maximum flow from s to t
func (d *Dinic) MaxFlow(s, t int) int {
   flow := 0
   const INF = int(1e9)
   for {
       d.bfs(s)
       if d.level[t] < 0 {
           break
       }
       for i := range d.iter {
           d.iter[i] = 0
       }
       for {
           f := d.dfs(s, t, INF)
           if f == 0 {
               break
           }
           flow += f
       }
   }
   return flow
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   A := make([]int, n)
   B := make([]int, n)
   sumA, sumB := 0, 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &A[i])
       sumA += A[i]
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &B[i])
       sumB += B[i]
   }
   if sumA != sumB {
       fmt.Fprintln(writer, "NO")
       return
   }
   // build graph: nodes 0..2n+1, s=0, t=2n+1
   s := 0
   t := 2*n + 1
   d := NewDinic(2*n + 2)
   // source to i, i to i+n, i+n to sink
   for i := 0; i < n; i++ {
       d.AddEdge(s, i+1, A[i])
       d.AddEdge(i+1, i+1+n, A[i])
       d.AddEdge(i+1+n, t, B[i])
   }
   const INF = int(1e9)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       d.AddEdge(x, y+n, INF)
       d.AddEdge(y, x+n, INF)
   }
   flow := d.MaxFlow(s, t)
   if flow != sumA {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   // retrieve flows
   ans := make([][]int, n)
   for i := range ans {
       ans[i] = make([]int, n)
   }
   for u := 1; u <= n; u++ {
       for _, e := range d.G[u] {
           if e.to >= n+1 && e.to <= n+n {
               j := e.to - (n + 1)
               revCap := d.G[e.to][e.rev].cap
               ans[u-1][j] = revCap
           }
       }
   }
   // output matrix
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           fmt.Fprint(writer, ans[i][j])
           if j+1 < n {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}
