package main

import (
   "fmt"
   "math"
)

func main() {
   var nn, rr int
   if _, err := fmt.Scanf("%d %d", &nn, &rr); err != nil {
       return
   }
   n := float64(nn)
   r := float64(rr)
   // angle for each central segment
   jiao := math.Pi / n
   // half of that angle
   jiao2 := jiao / 2
   zj := jiao
   // height from center to midpoint of side
   h := math.Sin(zj) * r
   // projection along center axis
   d := math.Cos(zj) * r
   // angle for triangle to compute db
   p := math.Pi/2 - zj - jiao2
   db := math.Tan(p) * h
   // total area
   s := (h*d - h*db) * n
   fmt.Printf("%.10f", s)
}
