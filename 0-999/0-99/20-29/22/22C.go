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

   var n, m, v int
   if _, err := fmt.Fscan(reader, &n, &m, &v); err != nil {
       return
   }
   // Check if possible to form a connected graph with m edges
   maxExtra := (n-1)*(n-2)/2 + 1
   if m < n-1 || m > maxExtra {
       writer.WriteString("-1\n")
       return
   }
   edges := 0
   var a, b int
   // First, connect v with all other nodes to ensure connectivity
   for i := 1; i <= n; i++ {
       if i == v {
           continue
       }
       a, b = i, v
       edges++
       if edges == n-1 {
           // reverse last edge representation (undirected)
           a, b = b, a
       }
       fmt.Fprintf(writer, "%d %d\n", a, b)
   }
   // Add remaining edges among other nodes until reach m
   for i := 1; i <= n && edges < m; i++ {
       if i == v || i == b {
           continue
       }
       for j := i + 1; j <= n && edges < m; j++ {
           if j == v || j == b {
               continue
           }
           fmt.Fprintf(writer, "%d %d\n", i, j)
           edges++
       }
   }
}
