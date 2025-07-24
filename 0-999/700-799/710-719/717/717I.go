package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type Pt3 struct{ x, y, z float64 }
type Pt2 struct{ x, y float64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   A3 := make([]Pt3, n)
   for i := 0; i < n; i++ {
       var x, y, z float64
       fmt.Fscan(in, &x, &y, &z)
       A3[i] = Pt3{x, y, z}
   }
   var m int
   fmt.Fscan(in, &m)
   B3 := make([]Pt3, m)
   for i := 0; i < m; i++ {
       var x, y, z float64
       fmt.Fscan(in, &x, &y, &z)
       B3[i] = Pt3{x, y, z}
   }
   // plane for A: normal = (A1-A0)x(A2-A0)
   v1 := sub3(A3[1], A3[0])
   v2 := sub3(A3[2], A3[0])
   normal := cross3(v1, v2)
   // basis e1,e2 on plane
   e1 := normalize3(v1)
   e2 := normalize3(cross3(normal, e1))
   // project A to 2D
   A2 := make([]Pt2, n)
   minx, miny := math.Inf(1), math.Inf(1)
   maxx, maxy := math.Inf(-1), math.Inf(-1)
   for i := 0; i < n; i++ {
       d := sub3(A3[i], A3[0])
       x := dot3(d, e1)
       y := dot3(d, e2)
       A2[i] = Pt2{x, y}
       if x < minx {
           minx = x
       }
       if x > maxx {
           maxx = x
       }
       if y < miny {
           miny = y
       }
       if y > maxy {
           maxy = y
       }
   }
   // count crossings
   var cntPosNeg, cntNegPos int
   d0 := dot3(normal, sub3(B3[0], A3[0]))
   for i := 0; i < m; i++ {
       j := (i + 1) % m
       p := B3[i]
       q := B3[j]
       d1 := d0
       d2 := dot3(normal, sub3(q, A3[0]))
       d0 = d2
       if d1*d2 >= 0 {
           continue
       }
       // intersection t
       t := d1 / (d1 - d2)
       xi := p.x + t*(q.x-p.x)
       yi := p.y + t*(q.y-p.y)
       zi := p.z + t*(q.z-p.z)
       // project to 2D
       dpt := Pt3{xi - A3[0].x, yi - A3[0].y, zi - A3[0].z}
       x2 := dot3(dpt, e1)
       y2 := dot3(dpt, e2)
       // bounding box test
       if x2 < minx || x2 > maxx || y2 < miny || y2 > maxy {
           // outside
       } else {
           if pointInPoly(x2, y2, A2) {
               if d1 > 0 {
                   cntPosNeg++
               } else {
                   cntNegPos++
               }
           }
       }
   }
   if cntPosNeg != cntNegPos {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

func sub3(a, b Pt3) Pt3 { return Pt3{a.x - b.x, a.y - b.y, a.z - b.z} }
func dot3(a, b Pt3) float64 { return a.x*b.x + a.y*b.y + a.z*b.z }
func cross3(a, b Pt3) Pt3 {
   return Pt3{a.y*b.z - a.z*b.y, a.z*b.x - a.x*b.z, a.x*b.y - a.y*b.x}
}
func normalize3(a Pt3) Pt3 {
   l := math.Sqrt(a.x*a.x + a.y*a.y + a.z*a.z)
   if l == 0 {
       return a
   }
   return Pt3{a.x / l, a.y / l, a.z / l}
}

// pointInPoly: ray casting to +x direction
func pointInPoly(x, y float64, P []Pt2) bool {
   inside := false
   n := len(P)
   j := n - 1
   for i := 0; i < n; j = i {
       yi := P[i].y
       yj := P[j].y
       xi := P[i].x
       xj := P[j].x
       // check edge j->i
       if (yi > y) != (yj > y) {
           // compute intersection x coord
           xint := (xj- xi)*(y- yi)/(yj- yi) + xi
           if xint > x {
               inside = !inside
           }
       }
       i++
   }
   return inside
}
