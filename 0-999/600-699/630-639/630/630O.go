package main

import (
   "fmt"
   "math"
)

func main() {
   var px, py, vx, vy float64
   if _, err := fmt.Scan(&px, &py, &vx, &vy); err != nil {
       return
   }
   var a, b, c, d int
   if _, err := fmt.Scan(&a, &b, &c, &d); err != nil {
       return
   }
   // normalize direction vector
   t := math.Sqrt(vx*vx + vy*vy)
   vx /= t
   vy /= t
   // perpendicular vector
   vx1 := vy
   vy1 := -vx
   // compute and output points
   fmt.Printf("%.10f %.10f\n", px+vx*float64(b), py+vy*float64(b))
   fmt.Printf("%.10f %.10f\n", px-vx1*float64(a)/2, py-vy1*float64(a)/2)
   fmt.Printf("%.10f %.10f\n", px-vx1*float64(c)/2, py-vy1*float64(c)/2)
   fmt.Printf("%.10f %.10f\n", px-vx1*float64(c)/2-vx*float64(d), py-vy1*float64(c)/2-vy*float64(d))
   fmt.Printf("%.10f %.10f\n", px+vx1*float64(c)/2-vx*float64(d), py+vy1*float64(c)/2-vy*float64(d))
   fmt.Printf("%.10f %.10f\n", px+vx1*float64(c)/2, py+vy1*float64(c)/2)
   fmt.Printf("%.10f %.10f\n", px+vx1*float64(a)/2, py+vy1*float64(a)/2)
}
