package main

import (
   "bufio"
   "fmt"
   "os"
)

// Dinic max flow
type edge struct { to, rev int; cap int64 }

type Dinic struct {
   n    int
   adj  [][]edge
   level []int
   ptr   []int
   q    []int
}

func NewDinic(n int) *Dinic {
   return &Dinic{n: n, adj: make([][]edge, n), level: make([]int, n), ptr: make([]int, n), q: make([]int, n)}
}

func (d *Dinic) AddEdge(u, v int, cap int64) {
   d.adj[u] = append(d.adj[u], edge{v, len(d.adj[v]), cap})
   d.adj[v] = append(d.adj[v], edge{u, len(d.adj[u]) - 1, 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level {
       d.level[i] = -1
   }
   dq := d.q
   qi, qj := 0, 0
   dq[qj] = s; d.level[s] = 0; qj++
   for qi < qj {
       u := dq[qi]; qi++
       for _, e := range d.adj[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               dq[qj] = e.to; qj++
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t int, f int64) int64 {
   if u == t || f == 0 {
       return f
   }
   for i := d.ptr[u]; i < len(d.adj[u]); i++ {
       e := &d.adj[u][i]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           pushed := d.dfs(e.to, t, min(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.adj[e.to][e.rev].cap += pushed
               return pushed
           }
       }
       d.ptr[u]++
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
   var flow int64
   for d.bfs(s, t) {
       for i := range d.ptr {
           d.ptr[i] = 0
       }
       for {
           pushed := d.dfs(s, t, 1<<60)
           if pushed == 0 {
               break
           }
           flow += pushed
       }
   }
   return flow
}

func min(a, b int64) int64 {
   if a < b { return a }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int64, n)
   b := make([]int64, n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i], &c[i])
   }
   l := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &l[i], &r[i])
   }
   u := make([]int, m)
   v := make([]int, m)
   dval := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &u[i], &v[i], &dval[i])
       u[i]--; v[i]--
   }
   // compute capacity per var and total nodes
   caps := make([]int, n)
   totalNodes := 0
   for i := 0; i < n; i++ {
       caps[i] = r[i] - l[i]
       totalNodes += caps[i]
   }
   // index mapping
   id := make([]int, n+1)
   id[0] = 0
   for i := 0; i < n; i++ {
       id[i+1] = id[i] + caps[i]
   }
   S := totalNodes
   T := totalNodes + 1
   dinic := NewDinic(totalNodes + 2)
   var base int64 = 0
   var sumPos int64 = 0
   inf := int64(1e18)
   // build nodes
   for i := 0; i < n; i++ {
       // base f(l[i])
       x0 := int64(l[i])
       base += a[i]*x0*x0 + b[i]*x0 + c[i]
       // deltas for k = 1..caps[i]
       for k := 1; k <= caps[i]; k++ {
           t := l[i] + k
           prev := t - 1
           w := a[i]*int64(t)*int64(t) + b[i]*int64(t) + c[i] - (a[i]*int64(prev)*int64(prev) + b[i]*int64(prev) + c[i])
           idx := id[i] + (k - 1)
           if w >= 0 {
               dinic.AddEdge(S, idx, w)
               sumPos += w
           } else {
               dinic.AddEdge(idx, T, -w)
           }
           // chain to previous
           if k > 1 {
               prevIdx := id[i] + (k - 2)
               dinic.AddEdge(idx, prevIdx, inf)
           }
       }
   }
   // constraints
   for ci := 0; ci < m; ci++ {
       ui, vi, dv := u[ci], v[ci], dval[ci]
       // xi <= xj + dv => xi >= l[i]+k => xj >= xi - dv
       for k := 1; k <= caps[ui]; k++ {
           xi := l[ui] + k
           need := xi - dv
           // find k' for j: need <= l[j] + k' => k' >= need - l[j]
           kp := need - l[vi]
           if kp > caps[vi] {
               // impossible, but feasible guaranteed
               continue
           }
           if kp <= 0 {
               continue
           }
           from := id[ui] + (k - 1)
           to := id[vi] + (kp - 1)
           dinic.AddEdge(from, to, inf)
       }
   }
   flow := dinic.MaxFlow(S, T)
   res := base + (sumPos - flow)
   fmt.Println(res)
}
