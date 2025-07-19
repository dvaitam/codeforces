package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

var xx, yy, vv float64
var x, y, v, rrGlobal float64
var tx, ty float64

// getPos computes the distance from origin to the point at fraction m along segment (x,y) to (tx,ty).
func getPos(m float64) float64 {
   qx := x + (tx-x)*m
   qy := y + (ty-y)*m
   return math.Hypot(qx, qy)
}

// getExt returns the arc components: (b, angle) for point (x0,y0) relative to circle of radius rrGlobal
func getExt(x0, y0 float64) (float64, float64) {
   c := math.Hypot(x0, y0)
   a := rrGlobal
   b := math.Sqrt(c*c - a*a)
   return b, math.Acos(a / c)
}

// check computes minimal path length at time t
func check(t float64) float64 {
   // ternary search on fraction along straight segment
   l2, r2 := 0.0, 1.0
   var m1, m2 float64
   for i := 0; i < 400; i++ {
       m := (r2 - l2) / 3.0
       m1 = l2 + m
       m2 = r2 - m
       if getPos(m1) <= getPos(m2) {
           r2 = m2
       } else {
           l2 = m1
       }
   }
   // if closest approach is outside circle, go straight
   if getPos(l2) >= rrGlobal {
       return math.Hypot(x-tx, y-ty)
   }
   // otherwise, combine tangents and arc
   spang := math.Abs(math.Atan2(y, x) - math.Atan2(ty, tx))
   if 2*math.Pi-spang < spang {
       spang = 2*math.Pi - spang
   }
   b1, ang1 := getExt(x, y)
   b2, ang2 := getExt(tx, ty)
   spang -= (ang1 + ang2)
   return b1 + b2 + spang*rrGlobal
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &xx, &yy, &vv)
   fmt.Fscan(reader, &x, &y, &v, &rrGlobal)
   // initial angle and circle radius
   stang := math.Atan2(yy, xx)
   // radius of circular path
   rr := math.Hypot(xx, yy)
   // binary search on time
   l, r := 0.0, 1e7
   for i := 0; i < 200; i++ {
       m := (l + r) / 2.0
       nang := stang + vv*m/rr
       tx = math.Cos(nang) * rr
       ty = math.Sin(nang) * rr
       if check(m) <= v*m {
           r = m
       } else {
           l = m
       }
   }
   fmt.Fprintf(writer, "%.16f", r)
}
