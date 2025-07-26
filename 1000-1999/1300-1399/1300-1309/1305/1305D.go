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

   var n int
   fmt.Fscan(reader, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // find initial leaves
   leaves := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if len(adj[i]) == 1 {
           leaves = append(leaves, i)
       }
   }
   // interactive elimination
   for {
       if len(leaves) == 1 {
           // found root
           fmt.Fprintf(writer, "! %d\n", leaves[0])
           writer.Flush()
           return
       }
       next := make([]int, 0, (len(leaves)+1)/2)
       // pair up leaves
       m := len(leaves)
       i := 0
       for i+1 < m {
           u, v := leaves[i], leaves[i+1]
           // query
           fmt.Fprintf(writer, "? %d %d\n", u, v)
           writer.Flush()
           var w int
           fmt.Fscan(reader, &w)
           if w == u || w == v {
               fmt.Fprintf(writer, "! %d\n", w)
               writer.Flush()
               return
           }
           next = append(next, w)
           i += 2
       }
       if i < m {
           // odd one out
           next = append(next, leaves[i])
       }
       leaves = next
   }
}
