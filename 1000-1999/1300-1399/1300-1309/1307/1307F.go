package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU structure
type DSU struct {
   p []int
}
func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := range p {
       p[i] = i
   }
   return &DSU{p}
}
func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}
func (d *DSU) Union(x, y int) {
   rx, ry := d.Find(x), d.Find(y)
   if rx != ry {
       d.p[ry] = rx
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k, r int
   fmt.Fscan(in, &n, &k, &r)
   adj := make([][]int, n)
   edges := make([][2]int, n-1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--; v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       edges[i] = [2]int{u, v}
   }
   // rest stops
   rest := make([]int, r)
   for i := 0; i < r; i++ {
       fmt.Fscan(in, &rest[i])
       rest[i]--
   }
   // multi-source BFS for dist to nearest rest
   const INF = 1<<60
   dist := make([]int, n)
   for i := range dist {
       dist[i] = -1
   }
   q := make([]int, 0, n)
   for _, u := range rest {
       dist[u] = 0
       q = append(q, u)
   }
   for qi := 0; qi < len(q); qi++ {
       u := q[qi]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               q = append(q, v)
           }
       }
   }
   // DSU on nodes for good edges
   dsu := NewDSU(n)
   for _, e := range edges {
       u, v := e[0], e[1]
       if dist[u]+dist[v]+1 <= k {
           dsu.Union(u, v)
       }
   }
   // map comps to ids
   compID := make([]int, n)
   idMap := make(map[int]int, n)
   cid := 0
   for i := 0; i < n; i++ {
       r := dsu.Find(i)
       if _, ok := idMap[r]; !ok {
           idMap[r] = cid
           cid++
       }
       compID[i] = idMap[r]
   }
   // build comp tree
   C := cid
   cadj := make([][]int, C)
   for _, e := range edges {
       u, v := e[0], e[1]
       cu, cv := compID[u], compID[v]
       if cu != cv {
           cadj[cu] = append(cadj[cu], cv)
           cadj[cv] = append(cadj[cv], cu)
       }
   }
   // LCA prep
   LOG := 1
   for (1<<LOG) <= C { LOG++ }
   up := make([][]int, LOG)
   depth := make([]int, C)
   for i := range up {
       up[i] = make([]int, C)
       for j := range up[i] {
           up[i][j] = -1
       }
   }
   var dfs func(u, p int)
   dfs = func(u, p int) {
       up[0][u] = p
       for _, v := range cadj[u] {
           if v == p { continue }
           depth[v] = depth[u] + 1
           dfs(v, u)
       }
   }
   for i := 0; i < C; i++ {
       if up[0][i] == -1 {
           dfs(i, -1)
       }
   }
   for j := 1; j < LOG; j++ {
       for i := 0; i < C; i++ {
           if up[j-1][i] < 0 {
               up[j][i] = -1
           } else {
               up[j][i] = up[j-1][ up[j-1][i] ]
           }
       }
   }
   lca := func(a, b int) int {
       if depth[a] < depth[b] {
           a, b = b, a
       }
       diff := depth[a] - depth[b]
       for j := 0; j < LOG; j++ {
           if diff>>j & 1 == 1 {
               a = up[j][a]
           }
       }
       if a == b {
           return a
       }
       for j := LOG-1; j >= 0; j-- {
           if up[j][a] != up[j][b] {
               a = up[j][a]
               b = up[j][b]
           }
       }
       return up[0][a]
   }
   // queries
   var vq int
   fmt.Fscan(in, &vq)
   for i := 0; i < vq; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--; b--
       ca, cb := compID[a], compID[b]
       // comp distance
       w := lca(ca, cb)
       distc := depth[ca] + depth[cb] - 2*depth[w]
       if distc <= k {
           out.WriteString("YES\n")
       } else {
           out.WriteString("NO\n")
       }
   }
}
