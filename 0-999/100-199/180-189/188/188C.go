package main

import (
   "fmt"
)

// gcd computes the greatest common divisor of a and b using Euclidean algorithm
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
   g := gcd(a, b)
   // Compute LCM = a/g * b to avoid overflow
   lcm := (a / g) * b
   fmt.Println(lcm)
}
