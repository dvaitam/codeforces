package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a connection to a node with a certain weight
type Edge struct {
   to int
   w  int64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]Edge, n+1)
   var u, v int
   var w int64
   var total int64
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &u, &v, &w)
       adj[u] = append(adj[u], Edge{to: v, w: w})
       adj[v] = append(adj[v], Edge{to: u, w: w})
       total += w
   }

   // Compute distances from city 1
   dist := make([]int64, n+1)
   visited := make([]bool, n+1)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   visited[1] = true
   var maxd int64
   for i := 0; i < len(queue); i++ {
       cur := queue[i]
       for _, e := range adj[cur] {
           if !visited[e.to] {
               visited[e.to] = true
               dist[e.to] = dist[cur] + e.w
               if dist[e.to] > maxd {
                   maxd = dist[e.to]
               }
               queue = append(queue, e.to)
           }
       }
   }

   // Minimal travel is twice the sum minus the farthest distance from city 1
   ans := total*2 - maxd
   fmt.Fprint(writer, ans)
}
