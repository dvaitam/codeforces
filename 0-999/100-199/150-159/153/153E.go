package main

import (
   "fmt"
   "math"
)

func main() {
   var n int
   // Read number of points
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   pts := make([][2]int, n)
   // Read points (x, y) each on separate lines
   for i := 0; i < n; i++ {
       fmt.Scan(&pts[i][0])
       fmt.Scan(&pts[i][1])
   }
   // Compute maximum squared distance
   maxd2 := 0.0
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           dx := float64(pts[i][0] - pts[j][0])
           dy := float64(pts[i][1] - pts[j][1])
           d2 := dx*dx + dy*dy
           if d2 > maxd2 {
               maxd2 = d2
           }
       }
   }
   // Output the maximum distance
   fmt.Printf("%.10f\n", math.Sqrt(maxd2))
}
