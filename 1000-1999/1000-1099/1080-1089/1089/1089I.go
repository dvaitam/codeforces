package main

import (
   "fmt"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   g := gcd(a, b)
   // compute lcm as a/g * b to avoid overflow
   l := a/g * b
   fmt.Println(l)
}
