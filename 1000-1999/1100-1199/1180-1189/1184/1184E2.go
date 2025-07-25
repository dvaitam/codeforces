package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const LOG = 17

// Edge represents a graph edge
type Edge struct {
   u, v, w, idx int
}

var parentDSU []int
var rankDSU []int

func find(u int) int {
   if parentDSU[u] != u {
       parentDSU[u] = find(parentDSU[u])
   }
   return parentDSU[u]
}

func union(u, v int) bool {
   ru, rv := find(u), find(v)
   if ru == rv {
       return false
   }
   if rankDSU[ru] < rankDSU[rv] {
       parentDSU[ru] = rv
   } else if rankDSU[ru] > rankDSU[rv] {
       parentDSU[rv] = ru
   } else {
       parentDSU[rv] = ru
       rankDSU[ru]++
   }
   return true
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       edges[i] = Edge{u, v, w, i}
   }

   // Build MST with Kruskal
   parentDSU = make([]int, n+1)
   rankDSU = make([]int, n+1)
   for i := 1; i <= n; i++ {
       parentDSU[i] = i
   }
   sortedEdges := make([]Edge, m)
   copy(sortedEdges, edges)
   sort.Slice(sortedEdges, func(i, j int) bool {
       return sortedEdges[i].w < sortedEdges[j].w
   })
   inMST := make([]bool, m)
   adj := make([][]struct{to, w int}, n+1)
   cnt := 0
   for _, e := range sortedEdges {
       if union(e.u, e.v) {
           inMST[e.idx] = true
           adj[e.u] = append(adj[e.u], struct{to, w int}{e.v, e.w})
           adj[e.v] = append(adj[e.v], struct{to, w int}{e.u, e.w})
           cnt++
           if cnt == n-1 {
               break
           }
       }
   }

   // Prepare LCA structures
   depth := make([]int, n+1)
   par := make([][LOG+1]int, n+1)
   maxW := make([][LOG+1]int, n+1)
   // BFS from node 1
   queue := make([]int, 0, n)
   depth[1] = 1
   queue = append(queue, 1)
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       for _, e := range adj[u] {
           v := e.to
           if depth[v] == 0 {
               depth[v] = depth[u] + 1
               par[v][0] = u
               maxW[v][0] = e.w
               queue = append(queue, v)
           }
       }
   }
   for k := 1; k <= LOG; k++ {
       for v := 1; v <= n; v++ {
           p := par[v][k-1]
           par[v][k] = par[p][k-1]
           maxW[v][k] = max(maxW[v][k-1], maxW[p][k-1])
       }
   }

   // Answer queries for non-tree edges
   ans := make([]int, m)
   for _, e := range edges {
       if inMST[e.idx] {
           continue
       }
       u, v := e.u, e.v
       mx := 0
       if depth[u] < depth[v] {
           u, v = v, u
       }
       diff := depth[u] - depth[v]
       for k := 0; k <= LOG; k++ {
           if diff&(1<<k) != 0 {
               mx = max(mx, maxW[u][k])
               u = par[u][k]
           }
       }
       if u != v {
           for k := LOG; k >= 0; k-- {
               if par[u][k] != 0 && par[u][k] != par[v][k] {
                   mx = max(mx, maxW[u][k])
                   mx = max(mx, maxW[v][k])
                   u = par[u][k]
                   v = par[v][k]
               }
           }
           mx = max(mx, maxW[u][0])
           mx = max(mx, maxW[v][0])
       }
       ans[e.idx] = mx
   }
   // Output results in input order for non-MST edges
   for i := 0; i < m; i++ {
       if !inMST[i] {
           fmt.Fprintln(out, ans[i])
       }
   }
