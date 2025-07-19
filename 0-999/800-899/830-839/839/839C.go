package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   depth := make([]int, n+1)
   order := make([]int, 0, n)
   // BFS to get order and parent, depth
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   depth[1] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           queue = append(queue, v)
       }
   }
   dp := make([]float64, n+1)
   // process in reverse order for DP
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       sum := 0.0
       cnt := 0
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           sum += dp[v]
           cnt++
       }
       if cnt > 0 {
           dp[u] = sum / float64(cnt)
       } else {
           dp[u] = float64(depth[u])
       }
   }
   // print result for root
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%.10f\n", dp[1])
}
