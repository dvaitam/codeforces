package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents an undirected edge with an identifier
type Edge struct { to, id int }

var (
   n, m    int
   g        [][]Edge
   disc, low []int
   timer     int
   isBridge  []bool
)

// dfs finds bridges using Tarjan's algorithm
func dfs(u, parentEdgeId int) {
   timer++
   disc[u] = timer
   low[u] = timer
   for _, e := range g[u] {
       v, id := e.to, e.id
       if id == parentEdgeId {
           continue
       }
       if disc[v] == 0 {
           dfs(v, id)
           if low[v] < low[u] {
               low[u] = low[v]
           }
           if low[v] > disc[u] {
               isBridge[id] = true
           }
       } else {
           if disc[v] < low[u] {
               low[u] = disc[v]
           }
       }
   }
}

var (
   comp      []int
   compCount int
)

// dfs2 labels 2-edge-connected components, skipping bridges
func dfs2(u int) {
   comp[u] = compCount
   for _, e := range g[u] {
       if isBridge[e.id] {
           continue
       }
       v := e.to
       if comp[v] == 0 {
           dfs2(v)
       }
   }
}

var (
   tree  [][]int
   up    [][]int
   depth []int
   LOG   int
)

// dfs3 preprocesses ancestors for LCA
func dfs3(u, p int) {
   up[u][0] = p
   for i := 1; i < LOG; i++ {
       up[u][i] = up[up[u][i-1]][i-1]
   }
   for _, v := range tree[u] {
       if v == p {
           continue
       }
       depth[v] = depth[u] + 1
       dfs3(v, u)
   }
}

// lca returns lowest common ancestor of u and v in the tree
func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   diff := depth[u] - depth[v]
   for i := 0; i < LOG; i++ {
       if diff&(1<<i) != 0 {
           u = up[u][i]
       }
   }
   if u == v {
       return u
   }
   for i := LOG - 1; i >= 0; i-- {
       if up[u][i] != up[v][i] {
           u = up[u][i]
           v = up[v][i]
       }
   }
   return up[u][0]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k int
   fmt.Fscan(reader, &n, &m)
   g = make([][]Edge, n+1)
   for i := 1; i <= m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       g[a] = append(g[a], Edge{b, i})
       g[b] = append(g[b], Edge{a, i})
   }
   disc = make([]int, n+1)
   low = make([]int, n+1)
   isBridge = make([]bool, m+1)
   // find all bridges
   dfs(1, -1)
   for i := 1; i <= n; i++ {
       if disc[i] == 0 {
           dfs(i, -1)
       }
   }
   // build components
   comp = make([]int, n+1)
   for i := 1; i <= n; i++ {
       if comp[i] == 0 {
           compCount++
           dfs2(i)
       }
   }
   // build bridge tree
   tree = make([][]int, compCount+1)
   for u := 1; u <= n; u++ {
       for _, e := range g[u] {
           v, id := e.to, e.id
           if isBridge[id] && comp[u] < comp[v] {
               cu, cv := comp[u], comp[v]
               tree[cu] = append(tree[cu], cv)
               tree[cv] = append(tree[cv], cu)
           }
       }
   }
   // prepare LCA
   for LOG = 1; (1 << LOG) <= compCount; LOG++ {}
   up = make([][]int, compCount+1)
   for i := range up {
       up[i] = make([]int, LOG)
   }
   depth = make([]int, compCount+1)
   // run dfs3 from each root (forest)
   for i := 1; i <= compCount; i++ {
       if up[i][0] == 0 {
           dfs3(i, i)
       }
   }
   // answer queries
   fmt.Fscan(reader, &k)
   for i := 0; i < k; i++ {
       var s, l int
       fmt.Fscan(reader, &s, &l)
       u, v := comp[s], comp[l]
       w := lca(u, v)
       dist := depth[u] + depth[v] - 2*depth[w]
       fmt.Fprintln(writer, dist)
   }
}
