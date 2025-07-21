package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   var count int64
   // Generate primitive Pythagorean triples via Euclid's formula
   for m := 2; m*m+1 <= n; m++ {
       for k := 1; k < m; k++ {
           // m and k have opposite parity and are coprime
           if (m-k)%2 == 1 && gcd(m, k) == 1 {
               c0 := m*m + k*k
               if c0 > n {
                   continue
               }
               // All multiples of the primitive triple with c = c0 * t <= n
               count += int64(n / c0)
           }
       }
   }
   fmt.Println(count)
}
