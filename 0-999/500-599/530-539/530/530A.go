package main

import (
   "fmt"
   "math"
)

func main() {
   var A, B, C float64
   if _, err := fmt.Scan(&A, &B, &C); err != nil {
       return
   }
   D := B*B - 4*A*C
   // discriminant is guaranteed non-negative
   sqrtD := math.Sqrt(D)
   x1 := (-B - sqrtD) / (2 * A)
   x2 := (-B + sqrtD) / (2 * A)
   if D == 0 {
       fmt.Printf("%.10f\n", x1)
   } else {
       if x1 > x2 {
           x1, x2 = x2, x1
       }
       fmt.Printf("%.10f %.10f\n", x1, x2)
   }
}
