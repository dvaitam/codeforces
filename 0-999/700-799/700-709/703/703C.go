package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func min(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b float64) float64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var w, v, u float64
   if _, err := fmt.Fscan(in, &n, &w, &v, &u); err != nil {
       return
   }
   pts := make([]struct{ x, y float64 }, n)
   for i := 0; i < n; i++ {
       var xi, yi float64
       fmt.Fscan(in, &xi, &yi)
       pts[i].x = xi
       pts[i].y = yi
   }
   direct := true
   ratio := v / u
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       xi, yi := pts[i].x, pts[i].y
       xj, yj := pts[j].x, pts[j].y
       if yi == yj {
           continue
       }
       // line: x(y) = xi + (xj-xi)*(y-yi)/(yj-yi) = A*(y-yi) + xi
       A := (xj - xi) / (yj - yi)
       // solve x(y) = ratio * y --> A*(y-yi) + xi = ratio*y
       // (A - ratio)*y + (xi - A*yi) = 0
       denom := A - ratio
       if math.Abs(denom) < 1e-12 {
           continue
       }
       y := (A*yi - xi) / denom
       // collision only if strictly inside segment
       if y > min(yi, yj) && y < max(yi, yj) {
           direct = false
           break
       }
   }
   tPed := w / u
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if direct {
       fmt.Fprintf(writer, "%.10f", tPed)
   } else {
       maxT := -1e300
       for i := 0; i < n; i++ {
           t := pts[i].x / v
           if t > maxT {
               maxT = t
           }
       }
       fmt.Fprintf(writer, "%.10f", maxT+tPed)
   }
}
