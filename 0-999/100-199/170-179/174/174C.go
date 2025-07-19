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
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   var (
       starts []int
       ends   []int
       s      int
   )
   // process heights, add sentinel zero at a[n+1]
   for i := 1; i <= n+1; i++ {
       x := a[i]
       // start segments when height increases
       for s < x {
           starts = append(starts, i)
           s++
       }
       // end segments when height decreases
       for s > x {
           ends = append(ends, i-1)
           s--
       }
   }

   m := len(starts)
   fmt.Fprintln(writer, m)
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, starts[i], ends[i])
   }
}
