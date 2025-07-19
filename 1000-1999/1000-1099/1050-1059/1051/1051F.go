package main

import (
   "bufio"
   "fmt"
   "os"
   "math"
)

const LG = 20
const INF = math.MaxInt64 / 4

type Edge struct {
   to int
   w  int64
}
type BEdge struct {
   u, v int
   w    int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   // DSU for tree
   dsu := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dsu[i] = i
   }
   find := func(x int) int {
       for dsu[x] != x {
           dsu[x] = dsu[dsu[x]]
           x = dsu[x]
       }
       return x
   }
   union := func(a, b int) {
       ra, rb := find(a), find(b)
       if ra != rb {
           dsu[ra] = rb
       }
   }

   tree := make([][]Edge, n+1)
   graph := make([][]Edge, n+1)
   bad := make([]BEdge, 0)

   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(in, &u, &v, &w)
       if find(u) != find(v) {
           union(u, v)
           tree[u] = append(tree[u], Edge{v, w})
           tree[v] = append(tree[v], Edge{u, w})
           graph[u] = append(graph[u], Edge{v, w})
           graph[v] = append(graph[v], Edge{u, w})
       } else {
           bad = append(bad, BEdge{u, v, w})
       }
   }
   // LCA prep
   parent := make([][]int, LG)
   for j := 0; j < LG; j++ {
       parent[j] = make([]int, n+1)
   }
   depth := make([]int, n+1)
   distRoot := make([]int64, n+1)
   // iterative DFS
   stack := []int{1}
   parent[0][1] = 0
   depth[1] = 0
   distRoot[1] = 0
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, e := range tree[u] {
           v := e.to
           if v == parent[0][u] {
               continue
           }
           parent[0][v] = u
           depth[v] = depth[u] + 1
           distRoot[v] = distRoot[u] + e.w
           stack = append(stack, v)
       }
   }
   for j := 1; j < LG; j++ {
       for i := 1; i <= n; i++ {
           pp := parent[j-1][i]
           if pp != 0 {
               parent[j][i] = parent[j-1][pp]
           }
       }
   }
   // augment graph with bad edges
   for _, e := range bad {
       graph[e.u] = append(graph[e.u], Edge{e.v, e.w})
       graph[e.v] = append(graph[e.v], Edge{e.u, e.w})
   }
   // SPFA distances per bad edge endpoints
   bsz := len(bad)
   dists := make([][]int64, 2*bsz)
   inQ := make([]bool, n+1)
   for i := 0; i < 2*bsz; i++ {
       dists[i] = make([]int64, n+1)
   }
   for i, e := range bad {
       spfa := func(s int, d []int64) {
           for k := 1; k <= n; k++ {
               d[k] = INF
               inQ[k] = false
           }
           que := make([]int, 0, n)
           head := 0
           d[s] = 0
           que = append(que, s)
           inQ[s] = true
           for head < len(que) {
               u := que[head]
               head++
               inQ[u] = false
               for _, ee := range graph[u] {
                   v := ee.to
                   nd := d[u] + ee.w
                   if nd < d[v] {
                       d[v] = nd
                       if !inQ[v] {
                           inQ[v] = true
                           que = append(que, v)
                       }
                   }
               }
           }
       }
       spfa(e.u, dists[2*i])
       spfa(e.v, dists[2*i+1])
   }
   // queries
   var q int
   fmt.Fscan(in, &q)
   // lca function
   lca := func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       diff := depth[u] - depth[v]
       for j := 0; j < LG; j++ {
           if diff&(1<<j) != 0 {
               u = parent[j][u]
           }
       }
       if u == v {
           return u
       }
       for j := LG - 1; j >= 0; j-- {
           if parent[j][u] != 0 && parent[j][u] != parent[j][v] {
               u = parent[j][u]
               v = parent[j][v]
           }
       }
       return parent[0][u]
   }
   for qi := 0; qi < q; qi++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       // tree distance
       w := distRoot[u] + distRoot[v] - 2*distRoot[lca(u, v)]
       ans := w
       for i, be := range bad {
           d1 := dists[2*i][u] + be.w + dists[2*i+1][v]
           if d1 < ans {
               ans = d1
           }
           d2 := dists[2*i][v] + be.w + dists[2*i+1][u]
           if d2 < ans {
               ans = d2
           }
       }
       fmt.Fprintln(out, ans)
   }
}
