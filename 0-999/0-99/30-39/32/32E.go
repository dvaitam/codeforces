package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type Point struct { x, y float64 }

func sub(a, b Point) Point { return Point{a.x - b.x, a.y - b.y} }
func add(a, b Point) Point { return Point{a.x + b.x, a.y + b.y} }
func mul(a Point, k float64) Point { return Point{a.x * k, a.y * k} }
func dot(a, b Point) float64 { return a.x*b.x + a.y*b.y }
func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func abs(x float64) float64 { if x < 0 { return -x } ; return x }

// orientation: >0 if c is left of ab, <0 if right, 0 if collinear
func orient(a, b, c Point) float64 {
   return cross(sub(b, a), sub(c, a))
}

// check if point c is on segment ab
func onSegment(a, b, c Point) bool {
   if abs(orient(a, b, c)) > 1e-9 {
       return false
   }
   minx, maxx := math.Min(a.x, b.x), math.Max(a.x, b.x)
   miny, maxy := math.Min(a.y, b.y), math.Max(a.y, b.y)
   return c.x >= minx-1e-9 && c.x <= maxx+1e-9 && c.y >= miny-1e-9 && c.y <= maxy+1e-9
}

// check if segments ab and cd intersect (including endpoints)
func segIntersect(a, b, c, d Point) bool {
   o1 := orient(a, b, c)
   o2 := orient(a, b, d)
   o3 := orient(c, d, a)
   o4 := orient(c, d, b)
   if o1*o2 < 0 && o3*o4 < 0 {
       return true
   }
   if abs(o1) < 1e-9 && onSegment(a, b, c) { return true }
   if abs(o2) < 1e-9 && onSegment(a, b, d) { return true }
   if abs(o3) < 1e-9 && onSegment(c, d, a) { return true }
   if abs(o4) < 1e-9 && onSegment(c, d, b) { return true }
   return false
}

// reflect point p across line ab
func reflectPoint(p, a, b Point) Point {
   // project p onto line ab
   ap := sub(p, a)
   ab := sub(b, a)
   ab2 := dot(ab, ab)
   if ab2 == 0 {
       return p
   }
   t := dot(ap, ab) / ab2
   proj := add(a, mul(ab, t))
   // reflection
   return add(mul(proj, 2), Point{-p.x, -p.y})
}

// intersection of lines p + t*r and q + u*s; returns (intersect point, exists, s)
func lineIntersect(p, r, q, s Point) (Point, bool, float64) {
   rxs := cross(r, s)
   if abs(rxs) < 1e-9 {
       return Point{}, false, 0
   }
   qp := sub(q, p)
   t := cross(qp, s) / rxs
   u := cross(qp, r) / rxs
   ip := add(p, mul(r, t))
   return ip, true, u
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var xv, yv, xp, yp float64
   var xw1, yw1, xw2, yw2 float64
   var xm1, ym1, xm2, ym2 float64
   fmt.Fscan(in, &xv, &yv)
   fmt.Fscan(in, &xp, &yp)
   fmt.Fscan(in, &xw1, &yw1, &xw2, &yw2)
   fmt.Fscan(in, &xm1, &ym1, &xm2, &ym2)
   V := Point{xv, yv}
   P := Point{xp, yp}
   W1 := Point{xw1, yw1}
   W2 := Point{xw2, yw2}
   M1 := Point{xm1, ym1}
   M2 := Point{xm2, ym2}
   // direct vision
   if !segIntersect(V, P, W1, W2) && !segIntersect(V, P, M1, M2) {
       fmt.Println("YES")
       return
   }
   // reflection
   // check same side of mirror line
   oV := orient(M1, M2, V)
   oP := orient(M1, M2, P)
   if oV*oP <= 0 {
       fmt.Println("NO")
       return
   }
   // reflect P across line M1-M2
   // compute projection
   ap := sub(P, M1)
   ab := sub(M2, M1)
   ab2 := dot(ab, ab)
   t := dot(ap, ab) / ab2
   proj := add(M1, mul(ab, t))
   Pp := Point{2*proj.x - P.x, 2*proj.y - P.y}
   // find intersection R of line V->Pp with mirror line M1->M2
   rVec := sub(Pp, V)
   sVec := ab
   R, ok, u := lineIntersect(V, rVec, M1, sVec)
   if !ok || u < -1e-9 || u > 1+1e-9 {
       fmt.Println("NO")
       return
   }
   // R must be different from V and P
   // check segments V-R and R-P don't hit wall
   if segIntersect(V, R, W1, W2) || segIntersect(R, P, W1, W2) {
       fmt.Println("NO")
       return
   }
   // valid reflection
   fmt.Println("YES")
}
