package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   // read adjacency matrix
   dist := make([][]int64, n)
   for i := 0; i < n; i++ {
       dist[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           var x int64
           fmt.Fscan(reader, &x)
           dist[i][j] = x
       }
   }
   // deletion order
   x := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i])
       x[i]-- // zero index
   }
   used := make([]bool, n)
   ans := make([]int64, n)

   // process additions in reverse deletion order
   for idx := n - 1; idx >= 0; idx-- {
       k := x[idx]
       used[k] = true
       // update shortest paths via k
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               // try path i->k->j
               via := dist[i][k] + dist[k][j]
               if via < dist[i][j] {
                   dist[i][j] = via
               }
           }
       }
       // sum distances among used vertices
       var sum int64
       for i := 0; i < n; i++ {
           if !used[i] {
               continue
           }
           for j := 0; j < n; j++ {
               if !used[j] {
                   continue
               }
               sum += dist[i][j]
           }
       }
       ans[idx] = sum
   }
   // output answers
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
