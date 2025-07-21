package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b (non-negative).
func gcd(a, b int) int {
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   // Read degrees
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // Read coefficients of P(x), keep leading a0
   var a0 int
   for i := 0; i <= n; i++ {
       var ai int
       fmt.Fscan(reader, &ai)
       if i == 0 {
           a0 = ai
       }
   }
   // Read coefficients of Q(x), keep leading b0
   var b0 int
   for i := 0; i <= m; i++ {
       var bi int
       fmt.Fscan(reader, &bi)
       if i == 0 {
           b0 = bi
       }
   }
   // Compare degrees
   switch {
   case n > m:
       // Infinity or -Infinity depending on sign of a0/b0
       if (a0 > 0 && b0 > 0) || (a0 < 0 && b0 < 0) {
           fmt.Println("Infinity")
       } else {
           fmt.Println("-Infinity")
       }
   case n < m:
       // Degree numerator lower: limit = 0
       fmt.Println("0/1")
   default:
       // n == m: limit = a0 / b0, reduce fraction
       num := a0
       den := b0
       // ensure denominator positive
       if den < 0 {
           num = -num
           den = -den
       }
       g := gcd(num, den)
       num /= g
       den /= g
       fmt.Printf("%d/%d\n", num, den)
   }
}
