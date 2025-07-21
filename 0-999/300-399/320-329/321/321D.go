package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   x := (n + 1) / 2
   var sumAbs int64
   var cntNeg int
   var minAbs int64 = -1
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var v int64
           fmt.Scan(&v)
           if v < 0 {
               cntNeg++
               v = -v
           }
           sumAbs += v
           if minAbs < 0 || v < minAbs {
               minAbs = v
           }
       }
   }
   // If flip size x is odd, we can adjust parity freely
   // Else, parity of negative count is invariant, so if it's odd, one must remain
   if x%2 == 0 && cntNeg%2 != 0 {
       sumAbs -= 2 * minAbs
   }
   fmt.Println(sumAbs)
}
