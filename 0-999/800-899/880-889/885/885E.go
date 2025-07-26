package main

import (
   "fmt"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a % b
   }
   return a
}

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   fmt.Println(gcd(a, b))
}
