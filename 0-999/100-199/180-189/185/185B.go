package main

import (
   "fmt"
)

func main() {
   var s, a, b, c int
   if _, err := fmt.Scan(&s); err != nil {
       return
   }
   _, _ = fmt.Scan(&a, &b, &c)
   // all zero exponents: any point yields distance 1^3, output zeros
   if a == 0 && b == 0 && c == 0 {
       fmt.Println("0 0 0")
       return
   }
   // handle case c == 0 separately to avoid division by zero
   if c == 0 {
       x := float64(s) * float64(a) / float64(a+b)
       y := float64(s) * float64(b) / float64(a+b)
       fmt.Printf("%.20f %.20f 0\n", x, y)
       return
   }
   S := float64(s)
   ca := float64(a)
   cb := float64(b)
   cc := float64(c)
   // using Lagrange multipliers: x = S*a/(a+b+c), etc.
   y := (S * cc * cb) / (cc*cb + cc*ca + cc*cc)
   x := (S * ca - y * ca) / (cc + ca)
   z := S - x - y
   fmt.Printf("%.20f %.20f %.20f\n", x, y, z)
}
