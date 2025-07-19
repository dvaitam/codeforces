package main

import (
   "fmt"
   "math"
)

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   var n int
   fmt.Scan(&n)
   minTime := math.Inf(1)
   for i := 0; i < n; i++ {
       var x, y, v int
       fmt.Scan(&x, &y, &v)
       dx := float64(x - a)
       dy := float64(y - b)
       dist := math.Sqrt(dx*dx + dy*dy)
       t := dist / float64(v)
       if t < minTime {
           minTime = t
       }
   }
   fmt.Printf("%f\n", minTime)
}
