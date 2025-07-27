package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Edge represents a possible road between airport u and factory v taking d days
type Edge struct {
   u, v int
   d    int
}

// hopcroftKarp computes maximum bipartite matching for left nodes 1..n
// adj is adjacency list from left to right nodes (1..n)
func hopcroftKarp(adj [][]int, n int) int {
   const inf = int(1e9)
   pairU := make([]int, n+1)
   pairV := make([]int, n+1)
   dist := make([]int, n+1)

   // BFS layers the graph, returns true if there is an augmenting path
   bfs := func() bool {
       queue := make([]int, 0, n)
       for u := 1; u <= n; u++ {
           if pairU[u] == 0 {
               dist[u] = 0
               queue = append(queue, u)
           } else {
               dist[u] = inf
           }
       }
       dist[0] = inf
       for i := 0; i < len(queue); i++ {
           u := queue[i]
           if dist[u] < dist[0] {
               for _, v := range adj[u] {
                   if pairV[v] == 0 {
                       if dist[0] == inf {
                           dist[0] = dist[u] + 1
                       }
                   } else if dist[pairV[v]] == inf {
                       dist[pairV[v]] = dist[u] + 1
                       queue = append(queue, pairV[v])
                   }
               }
           }
       }
       return dist[0] != inf
   }

   // dfs tries to find augmenting path from u
   var dfs func(u int) bool
   dfs = func(u int) bool {
       if u != 0 {
           for _, v := range adj[u] {
               if pairV[v] == 0 || (dist[pairV[v]] == dist[u]+1 && dfs(pairV[v])) {
                   pairU[u] = v
                   pairV[v] = u
                   return true
               }
           }
           dist[u] = inf
           return false
       }
       return true
   }

   matching := 0
   for bfs() {
       for u := 1; u <= n; u++ {
           if pairU[u] == 0 && dfs(u) {
               matching++
           }
       }
   }
   return matching
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       fmt.Fprintln(out, -1)
       return
   }
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].d)
   }
   sort.Slice(edges, func(i, j int) bool { return edges[i].d < edges[j].d })
   // collect unique days
   ds := make([]int, 0, m)
   last := -1
   for _, e := range edges {
       if e.d != last {
           ds = append(ds, e.d)
           last = e.d
       }
   }
   // binary search on ds
   ans := -1
   lo, hi := 0, len(ds)-1
   for lo <= hi {
       mid := (lo + hi) / 2
       threshold := ds[mid]
       // build graph with edges <= threshold
       adj := make([][]int, n+1)
       for _, e := range edges {
           if e.d > threshold {
               break
           }
           adj[e.u] = append(adj[e.u], e.v)
       }
       if hopcroftKarp(adj, n) == n {
           ans = threshold
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   if ans < 0 {
       fmt.Fprintln(out, -1)
   } else {
       fmt.Fprintln(out, ans)
   }
}
