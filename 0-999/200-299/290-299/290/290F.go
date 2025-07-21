package main

import (
   "bufio"
   "fmt"
   "os"
)

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
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       if u >= 0 && u < n && v >= 0 && v < n {
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
   }

   total := 1 << n
   // dp[mask*n+v] = reachable
   dp := make([]bool, total*n)
   prev := make([]int8, total*n)
   for v := 0; v < n; v++ {
       mask := 1 << v
       dp[mask*n+v] = true
       prev[mask*n+v] = -1
   }
   fullMask := total - 1
   found := false
   endV := -1
   for mask := 1; mask < total; mask++ {
       for v := 0; v < n; v++ {
           if !dp[mask*n+v] {
               continue
           }
           if mask == fullMask {
               found = true
               endV = v
               break
           }
           for _, u := range adj[v] {
               if mask&(1<<u) == 0 {
                   m2 := mask | (1 << u)
                   idx2 := m2*n + u
                   if !dp[idx2] {
                       dp[idx2] = true
                       prev[idx2] = int8(v)
                   }
               }
           }
       }
       if found {
           break
       }
   }
   if !found {
       fmt.Fprintln(writer, "No")
       return
   }
   // reconstruct path
   path := make([]int, n)
   mask := fullMask
   v := endV
   for i := n - 1; i >= 0; i-- {
       path[i] = v + 1
       p := prev[mask*n+v]
       mask ^= 1 << v
       v = int(p)
   }
   fmt.Fprintln(writer, "Yes")
   for i, x := range path {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, x)
   }
   fmt.Fprintln(writer)
}
