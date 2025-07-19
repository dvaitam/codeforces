package main

import (
   "fmt"
   "math"
)

func main() {
   var n int
   var h float64
   if _, err := fmt.Scan(&n, &h); err != nil {
       return
   }
   if n < 2 {
       return
   }
   // Compute ratios for recursive calculation
   out := make([]float64, n-1)
   out[0] = math.Sqrt(2)
   for i := 1; i < n-1; i++ {
       inv := 1 / (out[i-1] * out[i-1])
       out[i] = math.Sqrt(2 - inv)
   }
   // Compute cut positions from apex
   out1 := make([]float64, n-1)
   out1[n-2] = h / out[n-2]
   for i := n - 3; i >= 0; i-- {
       out1[i] = out1[i+1] / out[i]
   }
   // Print results with required precision
   for i := 0; i < n-1; i++ {
       fmt.Printf("%.12f", out1[i])
       if i+1 < n-1 {
           fmt.Print(" ")
       }
   }
   fmt.Println()
}
