package main

import (
   "fmt"
   "math"
)

func main() {
   var w, h, a float64
   const pi = math.Pi
   for {
       if _, err := fmt.Scan(&w, &h, &a); err != nil {
           break
       }
       // ensure w >= h
       if w < h {
           w, h = h, w
       }
       // diagonal length
       diag := math.Sqrt(w*w + h*h)
       // convert angle to radians and normalize
       a = a / 180.0 * pi
       if a > pi/2 {
           a = pi - a
       }
       // maximum overlap angle
       maxAngle := math.Asin(h/diag) * 2.0
       halfA := a / 2.0
       var ans float64
       if a >= maxAngle && a <= pi-maxAngle {
           // case when intersection is a hexagon-like shape
           l := h / math.Sin(halfA)
           hh := l/2.0 * math.Tan(halfA)
           ans = l * hh
       } else {
           // case when intersection is composed of triangles
           l := h * math.Sin(halfA)
           ans = l*l / math.Tan(halfA)
           g := math.Cos(a) + math.Sin(a) + 1.0
           gg := math.Cos(a) - math.Sin(a) + 1.0
           x := 0.0
           if gg != 0 {
               x = (w - (w+h)/g*math.Sin(a)) / gg
           }
           y := 0.0
           if g != 0 {
               y = (w+h)/g - x
           }
           s1 := x * math.Sin(halfA) * x * math.Cos(halfA) * 2.0
           s2 := y * math.Sin(halfA) * y * math.Cos(halfA) * 2.0
           s3 := y * math.Cos(halfA) * 2.0 * x * math.Cos(halfA) * 2.0
           ans = s1 + s2 + s3
       }
       fmt.Printf("%.7f\n", ans)
   }
}
