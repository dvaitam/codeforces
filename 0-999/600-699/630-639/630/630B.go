package main

import (
   "fmt"
   "math"
)

func main() {
   var n, t int64
   if _, err := fmt.Scan(&n, &t); err != nil {
       return
   }
   const a = 1.000000011
   var ans float64
   if t == 0 {
       ans = float64(n)
   } else {
       ans = float64(n) * math.Pow(a, float64(t))
   }
   fmt.Printf("%.6f", ans)
}
