package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge for Dinic
type Edge struct {
   to, rev, cap int
}

// Dinic implements max flow
type Dinic struct {
   n        int
   graph    [][]Edge
   level    []int
   ptr      []int
   queue    []int
}

// NewDinic creates a Dinic with n nodes
func NewDinic(n int) *Dinic {
   return &Dinic{
       n:     n,
       graph: make([][]Edge, n),
       level: make([]int, n),
       ptr:   make([]int, n),
       queue: make([]int, n),
   }
}

// AddEdge adds edge u->v with capacity cap
func (d *Dinic) AddEdge(u, v, cap int) {
   d.graph[u] = append(d.graph[u], Edge{to: v, rev: len(d.graph[v]), cap: cap})
   d.graph[v] = append(d.graph[v], Edge{to: u, rev: len(d.graph[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level {
       d.level[i] = -1
   }
   qh, qt := 0, 0
   d.queue[qt] = s
   qt++
   d.level[s] = 0
   for qh < qt {
       u := d.queue[qh]
       qh++
       for _, e := range d.graph[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               d.queue[qt] = e.to
               qt++
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t, pushed int) int {
   if u == t || pushed == 0 {
       return pushed
   }
   for i := d.ptr[u]; i < len(d.graph[u]); i++ {
       e := &d.graph[u][i]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           tr := d.dfs(e.to, t, min(pushed, e.cap))
           if tr > 0 {
               e.cap -= tr
               d.graph[e.to][e.rev].cap += tr
               return tr
           }
       }
       d.ptr[u]++
   }
   return 0
}

// Flow computes max flow from s to t
func (d *Dinic) Flow(s, t int) int {
   flow := 0
   for d.bfs(s, t) {
       for i := range d.ptr {
           d.ptr[i] = 0
       }
       for {
           pushed := d.dfs(s, t, 1<<30)
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
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n, m, k int
       fmt.Fscan(reader, &n, &m, &k)
       a := make([]int, m+1)
       b := make([]int, m+1)
       deg := make([]int, n+1)
       for i := 1; i <= m; i++ {
           fmt.Fscan(reader, &a[i], &b[i])
           deg[a[i]]++
           deg[b[i]]++
       }
       R := make([]int, n+1)
       sumR := 0
       for v := 1; v <= n; v++ {
           need := deg[v] - k
           if need < 0 {
               need = 0
           }
           R[v] = need * 2
           sumR += R[v]
       }
       Nnodes := n + m + 2
       S := 0
       T := n + m + 1
       dinic := NewDinic(Nnodes)
       for v := 1; v <= n; v++ {
           if R[v] > 0 {
               dinic.AddEdge(S, v, R[v])
           }
       }
       for i := 1; i <= m; i++ {
           en := n + i
           dinic.AddEdge(a[i], en, 1)
           dinic.AddEdge(b[i], en, 1)
           dinic.AddEdge(en, T, 1)
       }
       flow := dinic.Flow(S, T)
       comp := make([]int, m+1)
       if flow != sumR {
           // no solution
           for i := 1; i <= m; i++ {
               if i > 1 {
                   writer.WriteByte(' ')
               }
               writer.WriteString("0")
           }
           writer.WriteByte('\n')
           continue
       }
       sel := make([]int, m+1)
       for i := 1; i <= m; i++ {
           en := n + i
           // check a[i]
           for _, e := range dinic.graph[a[i]] {
               if e.to == en && e.cap == 0 {
                   sel[i] = a[i]
                   break
               }
           }
           if sel[i] == 0 {
               for _, e := range dinic.graph[b[i]] {
                   if e.to == en && e.cap == 0 {
                       sel[i] = b[i]
                       break
                   }
               }
           }
       }
       lists := make([][]int, n+1)
       for i := 1; i <= m; i++ {
           if sel[i] > 0 {
               lists[sel[i]] = append(lists[sel[i]], i)
           }
       }
       cid := 1
       for v := 1; v <= n; v++ {
           lst := lists[v]
           for j := 0; j+1 < len(lst); j += 2 {
               i1, i2 := lst[j], lst[j+1]
               comp[i1] = cid
               comp[i2] = cid
               cid++
           }
       }
       for i := 1; i <= m; i++ {
           if comp[i] == 0 {
               comp[i] = cid
               cid++
           }
       }
       for i := 1; i <= m; i++ {
           if i > 1 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprint(comp[i]))
       }
       writer.WriteByte('\n')
   }
}
