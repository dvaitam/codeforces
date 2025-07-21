package main

import (
   "fmt"
)

func main() {
   var x, y, m int64
   if _, err := fmt.Scan(&x, &y, &m); err != nil {
       return
   }
   // If already m-perfect
   if x >= m || y >= m {
       fmt.Println(0)
       return
   }
   // Impossible if both non-positive
   if x <= 0 && y <= 0 {
       fmt.Println(-1)
       return
   }
   // Ensure x <= y
   if x > y {
       x, y = y, x
   }
   ops := int64(0)
   // If x is negative, bring it non-negative
   if x < 0 {
       // y > 0 guaranteed here
       k := (-x + y - 1) / y
       x += k * y
       ops += k
   }
   // Now both x >= 0, y > 0
   for x < m && y < m {
       // add smaller to larger
       if x > y {
           x, y = y, x
       }
       x += y
       ops++
   }
   fmt.Println(ops)
}
