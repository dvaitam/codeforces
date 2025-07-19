package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const eps = 1e-8

// Pt represents a point or vector in 2D.
type Pt struct { x, y float64 }

func sub(a, b Pt) Pt   { return Pt{a.x - b.x, a.y - b.y} }
func cross(a, b Pt) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b Pt) float64   { return a.x*b.x + a.y*b.y }
func lenPt(a Pt) float64     { return math.Hypot(a.x, a.y) }

// orientation: 0=colinear, 1=clockwise, 2=counterclockwise
func orientation(p, q, r Pt) int {
   val := (q.y-p.y)*(r.x-q.x) - (q.x-p.x)*(r.y-q.y)
   if math.Abs(val) < eps {
       return 0
   }
   if val > 0 {
       return 1
   }
   return 2
}

// onSegment: check if q lies on segment pr
func onSegment(p, q, r Pt) bool {
   if q.x <= math.Max(p.x, r.x)+eps && q.x+eps >= math.Min(p.x, r.x) &&
      q.y <= math.Max(p.y, r.y)+eps && q.y+eps >= math.Min(p.y, r.y) {
       return true
   }
   return false
}

// doIntersect: check if segment p1q1 and p2q2 intersect
func doIntersect(p1, q1, p2, q2 Pt) bool {
   o1 := orientation(p1, q1, p2)
   o2 := orientation(p1, q1, q2)
   o3 := orientation(p2, q2, p1)
   o4 := orientation(p2, q2, q1)
   if o1 != o2 && o3 != o4 {
       return true
   }
   if o1 == 0 && onSegment(p1, p2, q1) { return true }
   if o2 == 0 && onSegment(p1, q2, q1) { return true }
   if o3 == 0 && onSegment(p2, p1, q2) { return true }
   if o4 == 0 && onSegment(p2, q1, q2) { return true }
   return false
}

// isInside: ray-casting to check if p is inside polygon
func isInside(polygon []Pt, p Pt, W, H float64) bool {
   n := len(polygon)
   if n < 3 {
       return false
   }
   mx := math.Max(W, H)
   extreme := Pt{mx*2 + 123.456, mx*2 + 789.012}
   count := 0
   for i := 0; i < n; i++ {
       next := (i + 1) % n
       if doIntersect(polygon[i], polygon[next], p, extreme) {
           if orientation(polygon[i], p, polygon[next]) == 0 {
               if onSegment(polygon[i], p, polygon[next]) {
                   return true
               }
           } else {
               count++
           }
       }
   }
   return count&1 == 1
}

// calc: distance from point x to segment ab, or large value if outside projection
func calc(x, a, b Pt) float64 {
   if dot(sub(x, a), sub(b, a)) <= 0 {
       return 1e200
   }
   if dot(sub(x, b), sub(a, b)) <= 0 {
       return 1e200
   }
   d := lenPt(sub(a, b))
   if d <= eps {
       return 1e200
   }
   return math.Abs(cross(sub(x, a), sub(x, b))) / d
}

// check intersection or containment between polygon and circle center p, radius r
func check(poly []Pt, p Pt, r float64, W, H float64) bool {
   // vertex distance
   for _, q := range poly {
       if lenPt(sub(p, q)) <= r+eps {
           return true
       }
   }
   // edge distance
   n := len(poly)
   for i := 0; i < n; i++ {
       if calc(p, poly[i], poly[(i+1)%n]) <= r+eps {
           return true
       }
   }
   // center inside polygon
   if isInside(poly, p, W, H) {
       return true
   }
   return false
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var W, H float64
   var n int
   fmt.Fscan(reader, &W, &H, &n)
   polys := make([][]Pt, n)
   areas := make([]float64, n)
   for i := 0; i < n; i++ {
       var cnt int
       fmt.Fscan(reader, &cnt)
       polys[i] = make([]Pt, cnt)
       for j := 0; j < cnt; j++ {
           fmt.Fscan(reader, &polys[i][j].x, &polys[i][j].y)
       }
       // compute area
       sum := 0.0
       for j := 0; j < cnt; j++ {
           k := (j + 1) % cnt
           sum += cross(polys[i][j], polys[i][k])
       }
       areas[i] = math.Abs(sum) / 2.0
   }
   var Q int
   fmt.Fscan(reader, &Q)
   for qi := 0; qi < Q; qi++ {
       var x, y, r float64
       fmt.Fscan(reader, &x, &y, &r)
       p := Pt{x, y}
       total := 0.0
       var res []int
       for i := 0; i < n; i++ {
           if check(polys[i], p, r, W, H) {
               total += areas[i]
               res = append(res, i)
           }
       }
       // output
       fmt.Fprintf(writer, "%.15f %d", total, len(res))
       for _, idx := range res {
           fmt.Fprintf(writer, " %d", idx)
       }
       writer.WriteByte('\n')
   }
}
