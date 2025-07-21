package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   degree := make([]int, n+1)
   for i := 0; i < n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
       degree[x]++
       degree[y]++
   }
   // find cycle nodes by leaf pruning
   inCycle := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       inCycle[i] = true
   }
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if degree[i] == 1 {
           queue = append(queue, i)
       }
   }
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       inCycle[u] = false
       for _, v := range adj[u] {
           if inCycle[v] {
               degree[v]--
               if degree[v] == 1 {
                   queue = append(queue, v)
               }
           }
       }
   }
   // BFS from cycle nodes to compute distances
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = -1
   }
   bfs := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if inCycle[i] {
           dist[i] = 0
           bfs = append(bfs, i)
       }
   }
   for head := 0; head < len(bfs); head++ {
       u := bfs[head]
       for _, v := range adj[u] {
           if dist[v] == -1 {
               dist[v] = dist[u] + 1
               bfs = append(bfs, v)
           }
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 1; i <= n; i++ {
       if i > 1 {
           w.WriteByte(' ')
       }
       fmt.Fprint(w, dist[i])
   }
   w.WriteByte('\n')
}
