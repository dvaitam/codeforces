package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to, w int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]edge, n+1)
   for i := 0; i < n-1; i++ {
       var u, v, w int
       fmt.Fscan(reader, &u, &v, &w)
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   // First BFS from node 1 to find farthest node u
   u, _ := bfs(1, adj, n)
   // Second BFS from u to find diameter length
   _, diameter := bfs(u, adj, n)
   fmt.Fprintln(writer, diameter)
}

// bfs returns the farthest node and its distance from start
func bfs(start int, adj [][]edge, n int) (farthestNode, farthestDist int) {
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = -1
   }
   queue := make([]int, 0, n)
   queue = append(queue, start)
   dist[start] = 0
   farthestNode = start
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       for _, e := range adj[u] {
           if dist[e.to] < 0 {
               dist[e.to] = dist[u] + e.w
               queue = append(queue, e.to)
               if dist[e.to] > dist[farthestNode] {
                   farthestNode = e.to
               }
           }
       }
   }
   farthestDist = dist[farthestNode]
   return
}
