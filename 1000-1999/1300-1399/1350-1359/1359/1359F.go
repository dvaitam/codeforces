package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const eps = 1e-8

// Point represents a point in 2D.
type Point struct {
   x, y float64
}

// Segment from p to q.
type Segment struct {
   p, q Point
}

// orientation returns 0 if colinear, 1 if clockwise, 2 if counterclockwise
func orientation(a, b, c Point) int {
   val := (b.y-a.y)*(c.x-b.x) - (b.x-a.x)*(c.y-b.y)
   if math.Abs(val) < eps {
       return 0
   }
   if val > 0 {
       return 1
   }
   return 2
}

// onSegment checks if point r lies on segment pq
func onSegment(p, q, r Point) bool {
   return r.x <= math.Max(p.x, q.x)+eps && r.x >= math.Min(p.x, q.x)-eps &&
       r.y <= math.Max(p.y, q.y)+eps && r.y >= math.Min(p.y, q.y)-eps
}

// segmentsIntersect checks if segment s1 intersects s2
func segmentsIntersect(s1, s2 Segment) bool {
   p1, q1, p2, q2 := s1.p, s1.q, s2.p, s2.q
   o1 := orientation(p1, q1, p2)
   o2 := orientation(p1, q1, q2)
   o3 := orientation(p2, q2, p1)
   o4 := orientation(p2, q2, q1)
   if o1 != o2 && o3 != o4 {
       return true
   }
   if o1 == 0 && onSegment(p1, q1, p2) {
       return true
   }
   if o2 == 0 && onSegment(p1, q1, q2) {
       return true
   }
   if o3 == 0 && onSegment(p2, q2, p1) {
       return true
   }
   if o4 == 0 && onSegment(p2, q2, q1) {
       return true
   }
   return false
}

// check determines if any segments intersect at time t
func check(t float64, xs, ys, dxs, dys, ss []float64) bool {
   n := len(xs)
   segs := make([]Segment, n)
   for i := 0; i < n; i++ {
       x, y := xs[i], ys[i]
       dx, dy := dxs[i], dys[i]
       speed := ss[i]
       norm := math.Hypot(dx, dy)
       // endpoint at time t
       ex := x + dx/norm*t*speed
       ey := y + dy/norm*t*speed
       segs[i] = Segment{Point{x, y}, Point{ex, ey}}
   }
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if segmentsIntersect(segs[i], segs[j]) {
               return true
           }
       }
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]float64, n)
   ys := make([]float64, n)
   dxs := make([]float64, n)
   dys := make([]float64, n)
   ss := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i], &dxs[i], &dys[i], &ss[i])
   }
   l, r, ans := 0.0, 1e10, 0.0
   for (r-l) > eps && (r-l)/r > eps {
       mid := (l + r) / 2
       if check(mid, xs, ys, dxs, dys, ss) {
           r = mid
           ans = mid
       } else {
           l = mid
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   if ans < eps {
       fmt.Fprintln(writer, "No show :(")
   } else {
       fmt.Fprintf(writer, "%.8f\n", ans)
   }
}
