package main

import (
   "fmt"
   "math"
)

func main() {
   var r, h int
   if _, err := fmt.Scan(&r, &h); err != nil {
       return
   }
   // Each balloon is a sphere of radius 3.5, so diameter = 7
   const diam = 7.0
   // Depth layers count
   dz := r / 7
   if dz <= 0 {
       fmt.Println(0)
       return
   }
   // Number of columns along width (front view)
   cols := (2 * r) / 7
   if cols <= 0 {
       fmt.Println(0)
       return
   }
   // Precompute effective semicircle radius for centers
   R0 := float64(r) - diam/2.0
   // Total balloons in one depth layer
   total2D := 0
   for i := 0; i < cols; i++ {
       // x coordinate of center: start at -r + 3.5, step by 7
       x := -float64(r) + diam/2.0 + float64(i)*diam
       // If outside semicircle horizontally, skip
       sq := x*x
       rr := R0 * R0
       if sq > rr {
           continue
       }
       // Maximum y coordinate allowed by semicircle
       yMax := float64(h) + math.Sqrt(rr - sq)
       // k_max = max integer k >= 0: 3.5 + k*7 <= yMax
       // => k <= (yMax - 3.5) / 7
       kMax := int(math.Floor((yMax - diam/2.0) / diam + 1e-9))
       if kMax >= 0 {
           total2D += kMax + 1
       }
   }
   // Total balloons is 2D count times depth layers
   result := int64(total2D) * int64(dz)
   fmt.Println(result)
}
