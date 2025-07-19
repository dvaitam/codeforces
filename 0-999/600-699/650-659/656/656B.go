package main

import (
   "fmt"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   m := make([]int64, n)
   r := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&m[i])
   }
   for i := 0; i < n; i++ {
       fmt.Scan(&r[i])
   }
   // Compute LCM of all m[i]
   var lcm int64 = 1
   for i := 0; i < n; i++ {
       g := gcd(lcm, m[i])
       lcm = lcm / g * m[i]
   }
   // Count numbers in [0, lcm) satisfying any congruence
   var count int64
   for i := int64(0); i < lcm; i++ {
       for j := 0; j < n; j++ {
           if i%m[j] == r[j] {
               count++
               break
           }
       }
   }
   result := float64(count) / float64(lcm)
   fmt.Printf("%.10f\n", result)
}
