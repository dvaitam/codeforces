package main

import (
   "fmt"
)

func main() {
   var nOrig int64
   if _, err := fmt.Scan(&nOrig); err != nil {
       return
   }
   if nOrig == 1 || nOrig == 2 {
       fmt.Println(-1)
       return
   }
   n := nOrig
   rec := 0
   for n%2 == 0 {
       n /= 2
       rec++
   }
   if n == 1 {
       x, y := int64(3), int64(5)
       // scale by 2^(rec-2)
       for i := 0; i < rec-2; i++ {
           x *= 2
           y *= 2
       }
       fmt.Println(x, y)
       return
   }
   x := n / 2
   y := x + 1
   ans1 := 2 * x * y
   ans2 := x*x + y*y
   for i := 0; i < rec; i++ {
       ans1 *= 2
       ans2 *= 2
   }
   fmt.Println(ans1, ans2)
}
