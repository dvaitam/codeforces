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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       var sumAbs int64
       var countNeg, countZero int
       var minAbs int64 = 1<<62
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               var x int64
               fmt.Fscan(reader, &x)
               if x < 0 {
                   countNeg++
               }
               if x == 0 {
                   countZero++
               }
               if x < 0 {
                   x = -x
               }
               if x < minAbs {
                   minAbs = x
               }
               sumAbs += x
           }
       }
       if countNeg%2 != 0 && countZero == 0 {
           sumAbs -= 2 * minAbs
       }
       fmt.Fprintln(writer, sumAbs)
   }
}
