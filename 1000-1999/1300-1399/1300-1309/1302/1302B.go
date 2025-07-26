package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   indegree := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       indegree[v]++
   }
   // Topological sort
   queue := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if indegree[i] == 0 {
           queue = append(queue, i)
       }
   }
   order := make([]int, 0, n)
   for idx := 0; idx < len(queue); idx++ {
       v := queue[idx]
       order = append(order, v)
       for _, u := range adj[v] {
           indegree[u]--
           if indegree[u] == 0 {
               queue = append(queue, u)
           }
       }
   }
   // Prepare bitsets
   words := (n + 63) >> 6
   bs := make([][]uint64, n)
   for i := 0; i < n; i++ {
       bs[i] = make([]uint64, words)
   }
   // Build reachable sets in reverse topological order
   for idx := n - 1; idx >= 0; idx-- {
       v := order[idx]
       for _, u := range adj[v] {
           // mark direct edge
           bs[v][u>>6] |= 1 << uint(u&63)
           // merge reachable of child
           for k := 0; k < words; k++ {
               bs[v][k] |= bs[u][k]
           }
       }
   }
   // Calculate answer
   var total uint64
   for i := 0; i < n; i++ {
       var cnt uint64
       for k := 0; k < words; k++ {
           cnt += uint64(bits.OnesCount64(bs[i][k]))
       }
       total += cnt * cnt
   }
   // Output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, total)
}
