package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const eps = 1e-9

// Point represents a point in 2D
type Point struct {
   x, y float64
}

func (a Point) sub(b Point) Point { return Point{a.x - b.x, a.y - b.y} }
func (a Point) add(b Point) Point { return Point{a.x + b.x, a.y + b.y} }
func (a Point) mul(s float64) Point { return Point{a.x * s, a.y * s} }
func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b Point) float64 { return a.x*b.x + a.y*b.y }

// segmentIntersect returns (t, exists) where intersection point is p = p + r*t
func segmentIntersect(p, r, q, s Point) (float64, bool) {
   // Solve p + r*t = q + s*u
   denom := cross(r, s)
   if math.Abs(denom) < eps {
       return 0, false
   }
   qp := q.sub(p)
   t := cross(qp, s) / denom
   u := cross(qp, r) / denom
   if t > -eps && t < 1+eps && u > -eps && u < 1+eps {
       return t, true
   }
   return 0, false
}

// pointInTriangle returns true if p is inside or on boundary of triangle abc
func pointInTriangle(p, a, b, c Point) bool {
   // Barycentric via cross signs
   ab := b.sub(a)
   bc := c.sub(b)
   ca := a.sub(c)
   ap := p.sub(a)
   bp := p.sub(b)
   cp := p.sub(c)
   c1 := cross(ab, ap)
   c2 := cross(bc, bp)
   c3 := cross(ca, cp)
   // check if all same sign or zero
   if (c1 >= -eps && c2 >= -eps && c3 >= -eps) || (c1 <= eps && c2 <= eps && c3 <= eps) {
       return true
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   tris := make([][3]Point, n)
   for i := 0; i < n; i++ {
       var x1, y1, x2, y2, x3, y3 float64
       fmt.Fscan(in, &x1, &y1, &x2, &y2, &x3, &y3)
       tris[i][0] = Point{x1, y1}
       tris[i][1] = Point{x2, y2}
       tris[i][2] = Point{x3, y3}
   }
   total := 0.0
   // For each triangle and each edge
   for i := 0; i < n; i++ {
       for e := 0; e < 3; e++ {
           a := tris[i][e]
           b := tris[i][(e+1)%3]
           r := b.sub(a)
           // collect t values
           ts := []float64{0.0, 1.0}
           // intersections with other triangle edges
           for j := 0; j < n; j++ {
               if j == i {
                   continue
               }
               for f := 0; f < 3; f++ {
                   c := tris[j][f]
                   d := tris[j][(f+1)%3]
                   s := d.sub(c)
                   if t, ok := segmentIntersect(a, r, c, s); ok {
                       if t > eps && t < 1-eps {
                           ts = append(ts, t)
                       }
                   }
               }
           }
           // sort and unique ts
           sort.Float64s(ts)
           // remove duplicates
           uniqueTs := ts[:0]
           for _, t := range ts {
               if len(uniqueTs) == 0 || math.Abs(t-uniqueTs[len(uniqueTs)-1]) > 1e-8 {
                   uniqueTs = append(uniqueTs, t)
               }
           }
           // for each interval, check midpoint
           for k := 0; k+1 < len(uniqueTs); k++ {
               t0 := uniqueTs[k]
               t1 := uniqueTs[k+1]
               tm := (t0 + t1) * 0.5
               pm := a.add(r.mul(tm))
               covered := false
               for j := 0; j < n; j++ {
                   if j == i {
                       continue
                   }
                   if pointInTriangle(pm, tris[j][0], tris[j][1], tris[j][2]) {
                       covered = true
                       break
                   }
               }
               if !covered {
                   total += (t1 - t0) * math.Hypot(r.x, r.y)
               }
           }
       }
   }
   fmt.Printf("%.6f\n", total)
}
