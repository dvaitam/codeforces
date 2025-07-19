package main

import (
   "bufio"
   "fmt"
   "os"
)

// coverall computes for each edge f and c counts for cover-all paths
func coverall(n, m, s, t int, edges [][4]int) ([]int, []int) {
   E := make([][]struct{v, idx int}, n)
   RE := make([][]struct{v, idx int}, n)
   for _, e := range edges {
       u, v, g, idx := e[0], e[1], e[2], e[3]
       if g == 0 {
           continue
       }
       E[u] = append(E[u], struct{v, idx int}{v, idx})
       RE[v] = append(RE[v], struct{v, idx int}{u, idx})
   }
   par := make([][2]int, n)
   rpar := make([][2]int, n)
   for i := 0; i < n; i++ {
       par[i][0], par[i][1] = -1, -1
       rpar[i][0], rpar[i][1] = -1, -1
   }
   par[s][0], par[s][1] = -2, -2
   rpar[t][0], rpar[t][1] = -2, -2
   // DFS from s
   stack := []int{s}
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, ch := range E[u] {
           if par[ch.v][0] != -1 {
               continue
           }
           par[ch.v][0], par[ch.v][1] = u, ch.idx
           stack = append(stack, ch.v)
       }
   }
   // DFS from t on reverse
   stack = []int{t}
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, ch := range RE[u] {
           if rpar[ch.v][0] != -1 {
               continue
           }
           rpar[ch.v][0], rpar[ch.v][1] = u, ch.idx
           stack = append(stack, ch.v)
       }
   }
   f := make([]int, m)
   c := make([]int, m)
   for _, e := range edges {
       u, v, g, idx := e[0], e[1], e[2], e[3]
       if g == 0 || f[idx] > 0 {
           continue
       }
       // include this edge
       f[idx]++
       c[idx]++
       // go up from u
       w := u
       for par[w][0] >= 0 {
           ei := par[w][1]
           f[ei]++
           c[ei]++
           w = par[w][0]
       }
       // go up from v in reverse
       w = v
       for rpar[w][0] >= 0 {
           ei := rpar[w][1]
           f[ei]++
           c[ei]++
           w = rpar[w][0]
       }
   }
   return f, c
}

// Edge for flow graph
type Edge struct {
   to, rev int
   cap     int64
   assoc   int
}

// Dinic max flow
type Dinic struct {
   n     int
   g     [][]Edge
   level []int
   ptr   []int
   q     []int
}

func NewDinic(n int) *Dinic {
   return &Dinic{n: n, g: make([][]Edge, n), level: make([]int, n), ptr: make([]int, n), q: make([]int, n)}
}

func (d *Dinic) AddEdge(u, v int, cap int64, assoc int) {
   d.g[u] = append(d.g[u], Edge{to: v, rev: len(d.g[v]), cap: cap, assoc: assoc})
   d.g[v] = append(d.g[v], Edge{to: u, rev: len(d.g[u]) - 1, cap: 0, assoc: -1})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := 0; i < d.n; i++ {
       d.level[i] = -1
   }
   d.level[s] = 0
   qh, qt := 0, 0
   d.q[qt] = s; qt++
   for qh < qt && d.level[t] < 0 {
       u := d.q[qh]; qh++
       for _, e := range d.g[u] {
           if e.cap > 0 && d.level[e.to] < 0 {
               d.level[e.to] = d.level[u] + 1
               d.q[qt] = e.to; qt++
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t int, pushed int64) int64 {
   if pushed == 0 {
       return 0
   }
   if u == t {
       return pushed
   }
   for d.ptr[u] < len(d.g[u]) {
       e := &d.g[u][d.ptr[u]]
       if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
           tr := d.dfs(e.to, t, minInt64(pushed, e.cap))
           if tr > 0 {
               e.cap -= tr
               d.g[e.to][e.rev].cap += tr
               return tr
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

func minInt64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

// block computes second part
func block(n, m, s, t int, edges [][4]int) ([]int, []int) {
   const INF = int64(1e12)
   d := NewDinic(n)
   for _, e := range edges {
       u, v, g, idx := e[0], e[1], e[2], e[3]
       if g == 0 {
           d.AddEdge(u, v, INF, idx)
       } else {
           d.AddEdge(u, v, 1, idx)
           d.AddEdge(v, u, INF, -1)
       }
   }
   d.MaxFlow(s, t)
   // reachable
   mc := make([]bool, n)
   var dfs func(int)
   dfs = func(u int) {
       mc[u] = true
       for _, e := range d.g[u] {
           if e.cap > 0 && !mc[e.to] {
               dfs(e.to)
           }
       }
   }
   dfs(s)
   f := make([]int, m)
   c := make([]int, m)
   for u := 0; u < n; u++ {
       for _, e := range d.g[u] {
           if e.cap == 0 {
               continue
           }
           if mc[u] && !mc[e.to] {
               continue
           }
           if e.assoc >= 0 {
               c[e.assoc]++
           }
       }
   }
   return f, c
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, s, t int
   fmt.Fscan(reader, &n, &m, &s, &t)
   s--; t--
   edges := make([][4]int, m)
   for i := 0; i < m; i++ {
       var u, v, g int
       fmt.Fscan(reader, &u, &v, &g)
       u--; v--
       edges[i] = [4]int{u, v, g, i}
   }
   f, c := coverall(n, m, s, t, edges)
   f2, c2 := block(n, m, s, t, edges)
   for i := 0; i < m; i++ {
       f[i] += f2[i]
       c[i] += c2[i]
       if c[i] == 0 {
           f[i] = 0
           c[i] = 2
       }
   }
   ans := 0
   for i := 0; i < m; i++ {
       if c[i] == f[i] {
           ans++
       }
   }
   fmt.Fprintln(writer, ans)
   for i := 0; i < m; i++ {
       fmt.Fprintf(writer, "%d %d\n", f[i], c[i])
   }
}
