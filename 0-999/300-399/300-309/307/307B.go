package main

import (
   "fmt"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   if b == 0 {
       if a < 0 {
           return -a
       }
       return a
   }
   return gcd(b, a%b)
}

func main() {
   var a, b int64
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   fmt.Println(gcd(a, b))
}
