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
   var x, y int64
   fmt.Fscan(reader, &n, &m, &x, &y)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &b[j])
   }

   i, j := 0, 0
   type pair struct{ u, v int }
   var res []pair
   for i < n && j < m {
       low := a[i] - x
       high := a[i] + y
       if b[j] < low {
           j++
       } else if b[j] > high {
           i++
       } else {
           // match
           res = append(res, pair{u: i + 1, v: j + 1})
           i++
           j++
       }
   }

   // output
   fmt.Fprintln(writer, len(res))
   for _, p := range res {
       fmt.Fprintln(writer, p.u, p.v)
   }
}
