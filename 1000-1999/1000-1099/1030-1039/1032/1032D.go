package main

import (
   "fmt"
   "math"
   "os"
)

func rast(x1, y1, x2, y2 float64) float64 {
   dx := x2 - x1
   dy := y2 - y1
   return math.Sqrt(dx*dx + dy*dy)
}

func main() {
   var a, b, c, x1, y1, x2, y2 float64
   fmt.Fscan(os.Stdin, &a, &b, &c)
   fmt.Fscan(os.Stdin, &x1, &y1, &x2, &y2)
   // direct Manhattan distance
   ans := math.Abs(x1-x2) + math.Abs(y1-y2)

   // if both coefficients non-zero, consider routes via line
   if a != 0 && b != 0 {
       tmp1y := (-c - a*x1) / b
       tmp1x := (-c - b*y1) / a
       tmp2y := (-c - a*x2) / b
       tmp2x := (-c - b*y2) / a
       r1x := rast(x1, y1, tmp1x, y1)
       r1y := rast(x1, y1, x1, tmp1y)
       r2x := rast(x2, y2, tmp2x, y2)
       r2y := rast(x2, y2, x2, tmp2y)
       // try all combinations of perpendiculars
       d1 := r1x + r2x + rast(tmp1x, y1, tmp2x, y2)
       d2 := r1x + r2y + rast(tmp1x, y1, x2, tmp2y)
       d3 := r1y + r2x + rast(x1, tmp1y, tmp2x, y2)
       d4 := r1y + r2y + rast(x1, tmp1y, x2, tmp2y)
       best := math.Min(math.Min(d1, d2), math.Min(d3, d4))
       if best < ans {
           ans = best
       }
   } else if a == 0 {
       // horizontal line: only y projection
       tmp1y := (-c - a*x1) / b
       tmp2y := (-c - a*x2) / b
       r1y := rast(x1, y1, x1, tmp1y)
       r2y := rast(x2, y2, x2, tmp2y)
       best := r1y + r2y + rast(x1, tmp1y, x2, tmp2y)
       if best < ans {
           ans = best
       }
   } else if b == 0 {
       // vertical line: only x projection
       tmp1x := (-c - b*y1) / a
       tmp2x := (-c - b*y2) / a
       r1x := rast(x1, y1, tmp1x, y1)
       r2x := rast(x2, y2, tmp2x, y2)
       best := r1x + r2x + rast(tmp1x, y1, tmp2x, y2)
       if best < ans {
           ans = best
       }
   }
   fmt.Printf("%.12f\n", ans)
}
