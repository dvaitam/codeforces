package main

import (
   "fmt"
)

func main() {
   var b, d, s int64
   if _, err := fmt.Scan(&b, &d, &s); err != nil {
       return
   }
   bs := []int64{b, d, s}
   sum := b + d + s
   const inf = int64(1<<62 - 1)
   ans := inf
   for A := int64(0); A <= 3; A++ {
       for D := int64(0); D <= 3; D++ {
           // elimination on first and last day for each meal
           var N int64
           for i := int64(0); i < 3; i++ {
               elF := int64(0)
               if i < A {
                   elF = 1
               }
               elL := int64(0)
               if i >= D {
                   elL = 1
               }
               need := bs[i] + elF + elL
               if need > N {
                   N = need
               }
           }
           if N == 0 {
               continue
           }
           // if only one day, arrival must not be after departure
           if N == 1 && A > D {
               continue
           }
           misses := 3*N - sum - (A + (3 - D))
           if misses < ans {
               ans = misses
           }
       }
   }
   fmt.Println(ans)
