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
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // Create grid
       q := make([][]int64, n)
       for i := 0; i < n; i++ {
           q[i] = make([]int64, m)
       }
       // Fill grid according to pattern
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               switch {
               case i == 0 && j == 0:
                   q[i][j] = 0
               case i == 0:
                   // First row, 1-based j+1
                   if (j+1)%2 == 1 {
                       q[i][j] = q[i][j-1] + 3
                   } else {
                       q[i][j] = q[i][j-1] + 1
                   }
               case (i+1)%2 == 0:
                   // Even row in 1-based indexing
                   q[i][j] = q[i-1][j] + 2
               default:
                   // Odd row (beyond first) in 1-based indexing
                   q[i][j] = q[i-2][j] + (1 << 9)
               }
           }
       }
       // Output number of distinct values
       fmt.Fprintln(writer, int64(n)*int64(m))
       // Output grid
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               if j > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, q[i][j])
           }
           fmt.Fprintln(writer)
       }
   }
}
