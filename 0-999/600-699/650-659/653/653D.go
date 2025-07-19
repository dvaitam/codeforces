package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Edge represents a flow edge
type Edge struct {
   to   int
   cap  int64
   rev  int
}

var (
   n, m int
   x     int64
   edges []struct{u, v int; c int}
   g     [][]Edge
   vis   []bool
)

// addEdge adds a directed edge u->v with capacity cap
func addEdge(u, v int, cap int64) {
   g[u] = append(g[u], Edge{to: v, cap: cap, rev: len(g[v])})
   g[v] = append(g[v], Edge{to: u, cap: 0,   rev: len(g[u]) - 1})
}

// dfs finds an augmenting path and returns flow
func dfs(u, t int, f int64) int64 {
   if u == t {
       return f
   }
   vis[u] = true
   for i := range g[u] {
       e := &g[u][i]
       if !vis[e.to] && e.cap > 0 {
           // send flow
           minf := f
           if e.cap < f {
               minf = e.cap
           }
           if ret := dfs(e.to, t, minf); ret > 0 {
               e.cap -= ret
               g[e.to][e.rev].cap += ret
               return ret
           }
       }
   }
   return 0
}

// can checks if with scale mid we can send at least x flow
func can(mid float64) bool {
   // build graph
   g = make([][]Edge, n)
   for i := 0; i < m; i++ {
       u, v, c := edges[i].u, edges[i].v, edges[i].c
       // capacity = floor(c / mid)
       cap := int64(math.Floor(float64(c) / mid))
       if cap < 0 {
           cap = 0
       }
       addEdge(u, v, cap)
   }
   // max flow
   s, t := 0, n-1
   var flow int64
   const INF = int64(1) << 60
   for flow < x {
       vis = make([]bool, n)
       f := dfs(s, t, INF)
       if f == 0 {
           break
       }
       flow += f
   }
   return flow >= x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &m, &x)
   edges = make([]struct{u, v int; c int}, m)
   maxc := 0
   for i := 0; i < m; i++ {
       var u, v, c int
       fmt.Fscan(in, &u, &v, &c)
       u--
       v--
       edges[i] = struct{u, v int; c int}{u, v, c}
       if c > maxc {
           maxc = c
       }
   }
   // binary search on average cost
   l, u := 0.0, float64(maxc)
   for it := 0; it < 80; it++ {
       mid := (l + u) / 2
       if mid == 0 {
           l = mid
           continue
       }
       if can(mid) {
           l = mid
       } else {
           u = mid
       }
   }
   ans := l * float64(x)
   fmt.Printf("%.20f\n", ans)
}
