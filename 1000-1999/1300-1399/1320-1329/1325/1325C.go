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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // edges: [u, v, ans]
   edges := make([][3]int, n-1)
   inDeg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       fmt.Fscan(reader, &edges[i][0], &edges[i][1])
       edges[i][2] = -1
       inDeg[edges[i][0]]++
       inDeg[edges[i][1]]++
   }
   // find a node with degree >= 3
   node := 0
   for i := 1; i <= n; i++ {
       if inDeg[i] >= 3 {
           node = i
           break
       }
   }
   if node == 0 {
       // no node with degree >=3, assign in order
       for i := 0; i < n-1; i++ {
           fmt.Fprintln(writer, i)
       }
       return
   }
   // assign labels
   cnt := 0
   for i := 0; i < n-1; i++ {
       if edges[i][0] == node || edges[i][1] == node {
           edges[i][2] = cnt
           cnt++
       }
   }
   for i := 0; i < n-1; i++ {
       if edges[i][2] == -1 {
           edges[i][2] = cnt
           cnt++
       }
   }
   for i := 0; i < n-1; i++ {
       fmt.Fprintln(writer, edges[i][2])
   }
}
