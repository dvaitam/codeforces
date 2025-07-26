package main

import (
   "fmt"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   var x, y int
   if _, err := fmt.Scan(&x, &y); err != nil {
       return
   }
   fmt.Println(gcd(x, y))
}
