package main

import "fmt"

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   var x, y int
   // read two positive integers
   if _, err := fmt.Scan(&x, &y); err != nil {
       return
   }
   // compute and print gcd
   fmt.Println(gcd(x, y))
}
