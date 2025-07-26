package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func bfs(start, n int, adj [][]int) []int {
   const INF = 1e9
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   q := make([]int, 0, n)
   dist[start] = 0
   q = append(q, start)
   for i := 0; i < len(q); i++ {
       u := q[i]
       du := dist[u]
       for _, v := range adj[u] {
           if dist[v] > du+1 {
               dist[v] = du + 1
               q = append(q, v)
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(in, &n, &m, &k)
   special := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &special[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   d1 := bfs(1, n, adj)
   d2 := bfs(n, n, adj)
   // original shortest path length
   L0 := d1[n]
   // sort special fields by d1 - d2
   sort.Slice(special, func(i, j int) bool {
       return d1[special[i]]-d2[special[i]] < d1[special[j]]-d2[special[j]]
   })
   best := 0
   mxD1 := -1000000000
   for _, u := range special {
       if mxD1 > -1000000000 {
           // candidate path using new road
           cand := mxD1 + d2[u] + 1
           if cand > best {
               best = cand
           }
       }
       // update maximum d1
       if d1[u] > mxD1 {
           mxD1 = d1[u]
       }
   }
   if best > L0 {
       best = L0
   }
   // if best is zero, fallback to original distance
   if best == 0 {
       best = L0
   }
   fmt.Println(best)
}
