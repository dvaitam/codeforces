package main

import (
   "fmt"
   "math"
)

func main() {
   var x1, y1, x2, y2, v, t, vx, vy, wx, wy float64
   // Read input values
   fmt.Scan(&x1, &y1, &x2, &y2)
   fmt.Scan(&v, &t)
   fmt.Scan(&vx, &vy, &wx, &wy)

   // Function to check if intercept is possible within time m
   check := func(m float64) bool {
       var tx, ty float64
       if m <= t {
           tx = x1 + vx*m
           ty = y1 + vy*m
       } else {
           tx = x1 + vx*t + wx*(m-t)
           ty = y1 + vy*t + wy*(m-t)
       }
       // Distance between target and interceptor
       return math.Hypot(tx-x2, ty-y2) <= v*m
   }

   // Binary search for minimal time
   l, r := 0.0, 1e8
   for i := 0; i < 200; i++ {
       mid := (l + r) * 0.5
       if check(mid) {
           r = mid
       } else {
           l = mid
       }
   }

   // Output result with precision
   fmt.Printf("%.12f", r)
}
