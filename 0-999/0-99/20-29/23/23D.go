package main

import (
   "bufio"
   "fmt"
   "os"
)

// cp returns the cross product of vectors (x0,y0)->(x1,y1) and (x0,y0)->(x2,y2)
func cp(x0, y0, x1, y1, x2, y2 float64) float64 {
   return (x1-x0)*(y2-y0) - (x2-x0)*(y1-y0)
}

// check attempts to find a valid quadruple given three points in order
// returns ok and the four points if successful
func check(x1, y1, x2, y2, x3, y3 float64) (bool, [8]float64) {
   sqr := func(x float64) float64 { return x * x }
   a1 := (x2 - x1) * 2.0
   b1 := (y2 - y1) * 2.0
   c1 := sqr(2*x1 - x2) + sqr(2*y1 - y2) - sqr(x1) - sqr(y1)
   a2 := (x3 - 2*x2 + x1) * 2.0
   b2 := (y3 - 2*y2 + y1) * 2.0
   c2 := sqr(x1) + sqr(y1) - sqr(x3-2*x2+2*x1) - sqr(y3-2*y2+2*y1)
   // parallel or degenerate
   if a1*b2 == a2*b1 {
       return false, [8]float64{}
   }
   Y1 := (c2*a1 - c1*a2) / (b1*a2 - b2*a1)
   X1 := (c2*b1 - c1*b2) / (a1*b2 - a2*b1)
   X2 := 2*x1 - X1
   Y2 := 2*y1 - Y1
   X3 := 2*x2 - 2*x1 + X1
   Y3 := 2*y2 - 2*y1 + Y1
   X4 := 2*x3 - 2*x2 + 2*x1 - X1
   Y4 := 2*y3 - 2*y2 + 2*y1 - Y1
   v1 := cp(X1, Y1, X2, Y2, X3, Y3)
   v2 := cp(X2, Y2, X3, Y3, X4, Y4)
   v3 := cp(X3, Y3, X4, Y4, X1, Y1)
   v4 := cp(X4, Y4, X1, Y1, X2, Y2)
   if (v1 < 0 && v2 < 0 && v3 < 0 && v4 < 0) || (v1 > 0 && v2 > 0 && v3 > 0 && v4 > 0) {
       return true, [8]float64{X1, Y1, X2, Y2, X3, Y3, X4, Y4}
   }
   return false, [8]float64{}
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var x1, y1, x2, y2, x3, y3 float64
       fmt.Fscan(in, &x1, &y1, &x2, &y2, &x3, &y3)
       found := false
       var pts [8]float64
       // try three orders
       if ok, p := check(x1, y1, x2, y2, x3, y3); ok {
           found, pts = true, p
       } else if ok, p := check(x1, y1, x3, y3, x2, y2); ok {
           found, pts = true, p
       } else if ok, p := check(x2, y2, x1, y1, x3, y3); ok {
           found, pts = true, p
       }
       if found {
           fmt.Fprintln(out, "YES")
           fmt.Fprintf(out, "%.9f %.9f %.9f %.9f %.9f %.9f %.9f %.9f\n",
               pts[0], pts[1], pts[2], pts[3], pts[4], pts[5], pts[6], pts[7])
       } else {
           fmt.Fprintln(out, "NO")
           fmt.Fprintln(out)
       }
   }
}
