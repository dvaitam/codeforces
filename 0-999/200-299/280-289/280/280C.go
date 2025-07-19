package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Scan(&x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   // BFS from node 1
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   depth[1] = 1
   parent[1] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, v := range adj[u] {
           if v != parent[u] {
               parent[v] = u
               depth[v] = depth[u] + 1
               queue = append(queue, v)
           }
       }
   }
   var ans float64
   for i := 1; i <= n; i++ {
       ans += 1.0 / float64(depth[i])
   }
   // print with 10 decimal places
   fmt.Printf("%.10f", ans)
}
