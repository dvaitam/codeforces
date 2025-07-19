package main

import (
   "fmt"
   "math"
)

// clc computes contribution based on distance q and radius r
func clc(q, r float64) float64 {
   if q < 1.0 {
       return 0.0
   }
   q1 := math.Min(1.0, q)
   base := math.Sqrt(r*r + r*r)
   rs := base*q1 + q1*q1*r
   q -= q1
   // For remaining q > 1
   return rs + base*q*2.0 + q*q*r
}

func main() {
   var j, k int
   if _, err := fmt.Scan(&j, &k); err != nil {
       return
   }
   n := float64(j)
   r := float64(k)
   rs := 0.0
   // Precompute base for efficiency
   for i := 1; i <= j; i++ {
       // two direct distances
       rs += r + r
       // left side
       q1 := float64(i - 1)
       rs += r*q1 + clc(q1, r)
       // right side
       q2 := n - float64(i)
       rs += r*q2 + clc(q2, r)
   }
   // average over all pairs
   rs /= (n * n)
   fmt.Printf("%.12f\n", rs)
}
