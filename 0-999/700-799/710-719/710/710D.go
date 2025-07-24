package main

import (
   "bufio"
   "fmt"
   "os"
)

// extGCD returns g = gcd(a,b) and x,y such that a*x + b*y = g
func extGCD(a, b int64) (g, x, y int64) {
   if b == 0 {
       return a, 1, 0
   }
   g, x1, y1 := extGCD(b, a%b)
   return g, y1, x1 - (a/b)*y1
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var a1, b1, a2, b2, L, R int64
   if _, err := fmt.Fscan(in, &a1, &b1, &a2, &b2, &L, &R); err != nil {
       return
   }
   // Compute solution of x ≡ b1 mod a1, x ≡ b2 mod a2
   // Solve a1*k - a2*l = b2 - b1
   diff := b2 - b1
   g, xg, _ := extGCD(a1, a2)
   if diff%g != 0 {
       fmt.Println(0)
       return
   }
   // one solution: k0 = (diff/g * xg) mod (a2/g)
   m2 := a2 / g
   k0 := (diff/g) * xg
   // normalize k0 into [0, m2)
   k0 %= m2
   if k0 < 0 {
       k0 += m2
   }
   // base solution x0
   x0 := b1 + a1*k0
   // lcm step
   step := a1 / g * a2
   // compute lower bound
   low := L
   if b1 > low {
       low = b1
   }
   if b2 > low {
       low = b2
   }
   if low > R {
       fmt.Println(0)
       return
   }
   // find minimal t such that x0 + t*step >= low
   var tMin int64
   if x0 < low {
       delta := low - x0
       tMin = (delta + step - 1) / step
   } else {
       tMin = 0
   }
   // find maximal t such that x0 + t*step <= R
   delta2 := R - x0
   if delta2 < 0 {
       fmt.Println(0)
       return
   }
   tMax := delta2 / step
   if tMin > tMax {
       fmt.Println(0)
   } else {
       fmt.Println(tMax - tMin + 1)
   }
}
