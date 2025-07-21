package main

import (
   "fmt"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

// floorDiv performs floor division of a by b.
func floorDiv(a, b int64) int64 {
   if b < 0 {
       a = -a
       b = -b
   }
   if a >= 0 {
       return a / b
   }
   return -(( -a + b - 1) / b)
}

func main() {
   var n, m, x, y, a, b int64
   if _, err := fmt.Scan(&n, &m, &x, &y, &a, &b); err != nil {
       return
   }
   // reduce ratio a:b
   g := gcd(a, b)
   a /= g
   b /= g
   // maximum scaling factor k
   k := n / a
   if km := m / b; km < k {
       k = km
   }
   // rectangle width and height
   wk := a * k
   hk := b * k
   // feasible x1 interval [lx, rx]
   lx := x - wk
   if lx < 0 {
       lx = 0
   }
   rx := x
   if n - wk < rx {
       rx = n - wk
   }
   // feasible y1 interval [ly, ry]
   ly := y - hk
   if ly < 0 {
       ly = 0
   }
   ry := y
   if m - hk < ry {
       ry = m - hk
   }
   // choose x1 minimizing distance to center
   D := 2*x - wk
   t0 := floorDiv(D, 2)
   // candidate positions
   candX := []int64{t0, t0 + 1}
   var x1 int64 = lx
   var bestDx int64 = -1
   for _, c := range candX {
       xx := c
       if xx < lx {
           xx = lx
       }
       if xx > rx {
           xx = rx
       }
       diff := 2*xx + wk - 2*x
       if diff < 0 {
           diff = -diff
       }
       if bestDx == -1 || diff < bestDx || (diff == bestDx && xx < x1) {
           bestDx = diff
           x1 = xx
       }
   }
   // choose y1 minimizing distance
   D = 2*y - hk
   t0 = floorDiv(D, 2)
   candY := []int64{t0, t0 + 1}
   var y1 int64 = ly
   var bestDy int64 = -1
   for _, c := range candY {
       yy := c
       if yy < ly {
           yy = ly
       }
       if yy > ry {
           yy = ry
       }
       diff := 2*yy + hk - 2*y
       if diff < 0 {
           diff = -diff
       }
       if bestDy == -1 || diff < bestDy || (diff == bestDy && yy < y1) {
           bestDy = diff
           y1 = yy
       }
   }
   x2 := x1 + wk
   y2 := y1 + hk
   fmt.Printf("%d %d %d %d", x1, y1, x2, y2)
}
