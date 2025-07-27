package main

import (
   "fmt"
)

// gcd computes the greatest common divisor of a and b using the Euclidean algorithm.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   // ensure non-negative inputs
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   fmt.Println(gcd(a, b))
}
