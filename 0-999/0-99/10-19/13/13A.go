package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var A int
   if _, err := fmt.Fscan(reader, &A); err != nil {
       return
   }
   sum := 0
   for base := 2; base < A; base++ {
       t := A
       for t > 0 {
           sum += t % base
           t /= base
       }
   }
   numerator := sum
   denominator := A - 2
   g := gcd(numerator, denominator)
   numerator /= g
   denominator /= g
   fmt.Printf("%d/%d", numerator, denominator)
}

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   if a < 0 {
       return -a
   }
   return a
}
