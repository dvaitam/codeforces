package main

import (
   "bufio"
   "fmt"
   "os"
)

func bfs(src int, adj [][]int) []int {
   n := len(adj)
   const inf = 1 << 30
   dist := make([]int, n)
   for i := range dist {
       dist[i] = inf
   }
   dist[src] = 0
   queue := make([]int, 0, n)
   queue = append(queue, src)
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       du := dist[u]
       for _, v := range adj[u] {
           if dist[v] > du+1 {
               dist[v] = du + 1
               queue = append(queue, v)
           }
       }
   }
   return dist
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   var s1, t1, l1, s2, t2, l2 int
   fmt.Fscan(reader, &s1, &t1, &l1)
   fmt.Fscan(reader, &s2, &t2, &l2)
   s1--
   t1--
   s2--
   t2--

   // compute all-pairs shortest paths
   dist := make([][]int, n)
   for i := 0; i < n; i++ {
       dist[i] = bfs(i, adj)
   }
   // check feasibility
   if dist[s1][t1] > l1 || dist[s2][t2] > l2 {
       fmt.Fprintln(writer, -1)
       return
   }
   // initial answer: disjoint shortest
   ans := dist[s1][t1] + dist[s2][t2]
   // try overlapping paths
   for u := 0; u < n; u++ {
       for v := 0; v < n; v++ {
           d_uv := dist[u][v]
           // path1 via u-v
           d1 := dist[s1][u] + d_uv + dist[v][t1]
           if d1 > l1 {
               continue
           }
           // path2 via u-v
           d2 := dist[s2][u] + d_uv + dist[v][t2]
           if d2 <= l2 {
               total := d1 + d2 - d_uv
               ans = min(ans, total)
           }
           // path2 via v-u
           d2 = dist[s2][v] + d_uv + dist[u][t2]
           if d2 <= l2 {
               total := d1 + d2 - d_uv
               ans = min(ans, total)
           }
       }
   }
   // maximum roads to destroy = m - used roads
   result := m - ans
   fmt.Fprintln(writer, result)
}
