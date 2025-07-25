package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1000000000

type Edge struct {
   u, v, w, idx int
}

// DSU for Kruskal
type DSU struct { p, r []int }
func newDSU(n int) *DSU {
   p := make([]int, n+1)
   r := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
       r[i] = 0
   }
   return &DSU{p, r}
}
func (d *DSU) find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.find(d.p[x])
   }
   return d.p[x]
}
func (d *DSU) union(x, y int) bool {
   x = d.find(x)
   y = d.find(y)
   if x == y {
       return false
   }
   if d.r[x] < d.r[y] {
       x, y = y, x
   }
   d.p[y] = x
   if d.r[x] == d.r[y] {
       d.r[x]++
   }
   return true
}

var (
   n, m int
   edges []Edge
   adj [][]struct{ to, w, idx int }
   inMST []bool
   parent [][]int
   maxW [][]int
   depth, sz, heavy, head, pos, parEdgeIdx []int
   curPos int
   segLazy []int
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   edges = make([]Edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
       edges[i].idx = i
   }
   // Kruskal
   dsu := newDSU(n)
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })
   inMST = make([]bool, m)
   adj = make([][]struct{ to, w, idx int }, n+1)
   for _, e := range edges {
       if dsu.union(e.u, e.v) {
           inMST[e.idx] = true
           adj[e.u] = append(adj[e.u], struct{ to, w, idx int }{e.v, e.w, e.idx})
           adj[e.v] = append(adj[e.v], struct{ to, w, idx int }{e.u, e.w, e.idx})
       }
   }
   // LCA prep and HLD prep
   parent = make([][]int, 18)
   maxW = make([][]int, 18)
   for i := 0; i < 18; i++ {
       parent[i] = make([]int, n+1)
       maxW[i] = make([]int, n+1)
   }
   depth = make([]int, n+1)
   sz = make([]int, n+1)
   heavy = make([]int, n+1)
   head = make([]int, n+1)
   pos = make([]int, n+1)
   parEdgeIdx = make([]int, n+1)
   // initial dfs from 1
   dfs1(1, 0)
   curPos = 1
   dfs2(1, 1)
   buildLCA()
   // segment tree lazy init
   segLazy = make([]int, 4*(n+1))
   for i := range segLazy {
       segLazy[i] = INF
   }
   ans := make([]int, m)
   // process non-tree edges
   for _, e := range edges {
       if !inMST[e.idx] {
           wmax := getMax(e.u, e.v)
           ans[e.idx] = wmax
           hldUpdate(e.u, e.v, e.w)
       }
   }
   // answers for tree edges
   for v := 2; v <= n; v++ {
       idx := parEdgeIdx[v]
       val := hldGet(pos[v])
       if val == INF {
           val = INF
       }
       ans[idx] = val
   }
   // output
   for i := 0; i < m; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

func dfs1(u, p int) {
   sz[u] = 1
   heavy[u] = 0
   for _, e := range adj[u] {
       v := e.to
       if v == p {
           continue
       }
       parent[0][v] = u
       maxW[0][v] = e.w
       depth[v] = depth[u] + 1
       parEdgeIdx[v] = e.idx
       dfs1(v, u)
       if sz[v] > sz[heavy[u]] {
           heavy[u] = v
       }
       sz[u] += sz[v]
   }
}

func dfs2(u, h int) {
   head[u] = h
   pos[u] = curPos
   curPos++
   if heavy[u] != 0 {
       dfs2(heavy[u], h)
   }
   for _, e := range adj[u] {
       v := e.to
       if v == parent[0][u] || v == heavy[u] {
           continue
       }
       dfs2(v, v)
   }
}

func buildLCA() {
   for k := 1; k < 18; k++ {
       for v := 1; v <= n; v++ {
           p := parent[k-1][v]
           parent[k][v] = parent[k-1][p]
           maxW[k][v] = max(maxW[k-1][v], maxW[k-1][p])
       }
   }
}

func getMax(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   var res int
   d := depth[u] - depth[v]
   for k := 0; d > 0; k++ {
       if d&1 != 0 {
           res = max(res, maxW[k][u])
           u = parent[k][u]
       }
       d >>= 1
   }
   if u == v {
       return res
   }
   for k := 17; k >= 0; k-- {
       if parent[k][u] != parent[k][v] {
           res = max(res, maxW[k][u])
           res = max(res, maxW[k][v])
           u = parent[k][u]
           v = parent[k][v]
       }
   }
   res = max(res, maxW[0][u])
   res = max(res, maxW[0][v])
   return res
}

// HLD update: cap edges on path u-v with w
func hldUpdate(u, v, w int) {
   for head[u] != head[v] {
       if depth[head[u]] > depth[head[v]] {
           updateRange(pos[head[u]], pos[u], w, 1, 1, n)
           u = parent[0][head[u]]
       } else {
           updateRange(pos[head[v]], pos[v], w, 1, 1, n)
           v = parent[0][head[v]]
       }
   }
   if u == v {
       return
   }
   if depth[u] > depth[v] {
       updateRange(pos[v]+1, pos[u], w, 1, 1, n)
   } else {
       updateRange(pos[u]+1, pos[v], w, 1, 1, n)
   }
}

func updateRange(l, r, w, node, nl, nr int) {
   if l > nr || r < nl {
       return
   }
   if l <= nl && nr <= r {
       if segLazy[node] > w {
           segLazy[node] = w
       }
       return
   }
   mid := (nl + nr) >> 1
   left, right := node<<1, node<<1|1
   if segLazy[node] < INF {
       if segLazy[left] > segLazy[node] {
           segLazy[left] = segLazy[node]
       }
       if segLazy[right] > segLazy[node] {
           segLazy[right] = segLazy[node]
       }
       segLazy[node] = INF
   }
   updateRange(l, r, w, left, nl, mid)
   updateRange(l, r, w, right, mid+1, nr)
}

// get cap value at position p
func hldGet(p int) int {
   return getVal(p, 1, 1, n)
}

func getVal(posi, node, nl, nr int) int {
   if nl == nr {
       return segLazy[node]
   }
   mid := (nl + nr) >> 1
   left, right := node<<1, node<<1|1
   if segLazy[node] < INF {
       if segLazy[left] > segLazy[node] {
           segLazy[left] = segLazy[node]
       }
       if segLazy[right] > segLazy[node] {
           segLazy[right] = segLazy[node]
       }
       segLazy[node] = INF
   }
   if posi <= mid {
       return getVal(posi, left, nl, mid)
   }
   return getVal(posi, right, mid+1, nr)
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
