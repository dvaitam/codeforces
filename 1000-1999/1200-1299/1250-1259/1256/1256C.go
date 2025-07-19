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

   var n, m, d int
   if _, err := fmt.Fscan(reader, &n, &m, &d); err != nil {
       return
   }
   c := make([]int, m)
   sum := 0
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &c[i])
       sum += c[i]
   }
   remZeros := n - sum
   // Check feasibility: total zeros must fit in (m+1) gaps of size <= d
   if remZeros > (m+1)*d {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   // Distribute zeros greedily up to d per gap
   gaps := make([]int, m+1)
   for i := 0; i <= m; i++ {
       if remZeros > d {
           gaps[i] = d
           remZeros -= d
       } else {
           gaps[i] = remZeros
           remZeros = 0
       }
   }
   // Build and output sequence
   for i := 0; i <= m; i++ {
       // zeros in gap i
       for j := 0; j < gaps[i]; j++ {
           fmt.Fprint(writer, 0, " ")
       }
       // segment i (1-indexed) except after last gap
       if i < m {
           for j := 0; j < c[i]; j++ {
               fmt.Fprint(writer, i+1, " ")
           }
       }
   }
   fmt.Fprintln(writer)
}
