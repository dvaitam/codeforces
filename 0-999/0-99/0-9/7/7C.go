package main

import (
   "fmt"
)

// extGCD computes gcd(a, b) and finds x, y such that a*x + b*y = gcd(a, b)
func extGCD(a, b int64) (g, x, y int64) {
   if b == 0 {
       return a, 1, 0
   }
   g2, x1, y1 := extGCD(b, a%b)
   // x1 * (b) + y1 * (a%b) = g2
   // a%b = a - (a/b)*b
   // => x1*b + y1*(a - (a/b)*b) = g2
   // => y1*a + (x1 - y1*(a/b))*b = g2
   return g2, y1, x1 - (a/b)*y1
}

func main() {
   var A, B, C int64
   if _, err := fmt.Scan(&A, &B, &C); err != nil {
       return
   }
   // Solve A*x + B*y + C = 0 => A*x + B*y = -C
   g, x0, y0 := extGCD(abs(A), abs(B))
   // Adjust signs for original A, B
   if A < 0 {
       x0 = -x0
   }
   if B < 0 {
       y0 = -y0
   }
   // Check divisibility
   if (-C)%g != 0 {
       fmt.Println(-1)
       return
   }
   factor := (-C) / g
   x := x0 * factor
   y := y0 * factor
   // Output solution
   fmt.Println(x, y)
}

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}
