package main

import (
   "fmt"
)

func main() {
   var x, y, n int64
   if _, err := fmt.Scan(&x, &y, &n); err != nil {
       return
   }
   const INF = 1 << 62
   bestNum := INF
   var bestA, bestDen int64 = 0, 1
   for b := int64(1); b <= n; b++ {
       // a = floor(x*b/y)
       a := x * b / y
       // consider a and a+1
       for _, ai := range []int64{a, a + 1} {
           if ai < 0 {
               continue
           }
           // absolute difference numerator: |x*b - ai*y|
           diff := x*b - ai*y
           if diff < 0 {
               diff = -diff
           }
           // compare diff/b with bestNum/bestDen: cross-multiply
           if diff*bestDen < bestNum*b ||
              (diff*bestDen == bestNum*b && (b < bestDen || (b == bestDen && ai < bestA))) {
               bestNum = diff
               bestDen = b
               bestA = ai
           }
       }
   }
   fmt.Printf("%d/%d\n", bestA, bestDen)
}
