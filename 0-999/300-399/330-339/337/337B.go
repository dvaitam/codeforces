package main

import (
   "fmt"
)

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

func main() {
   var a, b, c, d int
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   // Determine limiting dimension: width or height
   // Compare a/c and b/d via cross-multiplication: a*d <= b*c => width limiting
   t1 := a * d
   t2 := b * c
   var p, q int
   if t1 <= t2 {
       // width limiting: movie width = a, height = a*d/c
       p = b*c - a*d
       q = b * c
   } else {
       // height limiting: movie height = b, width = b*c/d
       p = a*d - b*c
       q = a * d
   }
   // simplify fraction
   if p == 0 {
       fmt.Println("0/1")
       return
   }
   g := gcd(p, q)
   p /= g
   q /= g
   fmt.Printf("%d/%d", p, q)
}
