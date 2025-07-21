package main

import (
   "fmt"
   "math"
)

func main() {
   var k, xi, yi int
   if _, err := fmt.Scan(&k, &xi, &yi); err != nil {
       return
   }
   x := float64(xi)
   y := float64(yi)
   theta := float64(k) * math.Pi / 180.0
   c := math.Cos(theta)
   s := math.Sin(theta)
   x2 := x*c - y*s
   y2 := x*s + y*c
   fmt.Printf("%.10f %.10f\n", x2, y2)
}
