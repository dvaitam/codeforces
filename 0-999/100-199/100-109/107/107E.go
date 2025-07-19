package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const EPS = 1e-8

// PT represents a point or vector
type PT struct {
   x, y float64
}

// Sub returns the vector subtraction p - q
func (p PT) Sub(q PT) PT { return PT{p.x - q.x, p.y - q.y} }

// Cross returns the cross product p x q
func (p PT) Cross(q PT) float64 { return p.x*q.y - p.y*q.x }

// SG returns the sign of x with EPS tolerance
func SG(x float64) int {
   if x > EPS {
       return 1
   }
   if x < -EPS {
       return -1
   }
   return 0
}

// tri returns the signed area*2 of triangle p1,p2,p3
func tri(p1, p2, p3 PT) float64 {
   return p2.Sub(p1).Cross(p3.Sub(p1))
}

// segP returns parameter t in [0,1] for projection of p on segment p1->p2
func segP(p, p1, p2 PT) float64 {
   if SG(p1.x - p2.x) == 0 {
       return (p.y - p1.y) / (p2.y - p1.y)
   }
   return (p.x - p1.x) / (p2.x - p1.x)
}

// polyUnion computes the union area of given polygons
func polyUnion(polys [][]PT) float64 {
   n := len(polys)
   var sum float64
   // process each edge of each polygon
   for i := 0; i < n; i++ {
       p := polys[i]
       m := len(p)
       // append first point at end
       p2 := make([]PT, m+1)
       copy(p2, p)
       p2[m] = p[0]
       for ii := 0; ii < m; ii++ {
           a, b := p2[ii], p2[ii+1]
           // events on [0,1]
           type interval struct{ t float64; d int }
           ev := make([]interval, 0, 4)
           ev = append(ev, interval{0, 0}, interval{1, 0})
           // for each other polygon
           for j := 0; j < n; j++ {
               if j == i {
                   continue
               }
               q := polys[j]
               lq := len(q)
               // append first point at end
               q2 := make([]PT, lq+1)
               copy(q2, q)
               q2[lq] = q[0]
               for jj := 0; jj < lq; jj++ {
                   c, dpt := q2[jj], q2[jj+1]
                   ta := SG(tri(a, b, c))
                   tb := SG(tri(a, b, dpt))
                   if ta == 0 && tb == 0 {
                       // collinear
                       if q2[jj+1].Sub(q2[jj]).Cross(b.Sub(a)) > 0 && j < i {
                           t1 := segP(c, a, b)
                           t2 := segP(dpt, a, b)
                           ev = append(ev, interval{t1, 1}, interval{t2, -1})
                       }
                   } else if ta >= 0 && tb < 0 {
                       tc := tri(c, dpt, a)
                       td := tri(c, dpt, b)
                       ev = append(ev, interval{tc / (tc - td), 1})
                   } else if ta < 0 && tb >= 0 {
                       tc := tri(c, dpt, a)
                       td := tri(c, dpt, b)
                       ev = append(ev, interval{tc / (tc - td), -1})
                   }
               }
           }
           // sort events by t
           sort.Slice(ev, func(i, j int) bool { return ev[i].t < ev[j].t })
           // sweep intervals
           z := math.Min(math.Max(ev[0].t, 0.0), 1.0)
           dcnt := ev[0].d
           var covered float64
           for k := 1; k < len(ev); k++ {
               w := math.Min(math.Max(ev[k].t, 0.0), 1.0)
               if dcnt == 0 {
                   covered += w - z
               }
               dcnt += ev[k].d
               z = w
           }
           sum += a.Cross(b) * covered
       }
   }
   return sum / 2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   polygons := make([][]PT, n)
   var sumArea float64
   for i := 0; i < n; i++ {
       pts := make([]PT, 4)
       for j := 0; j < 4; j++ {
           fmt.Fscan(reader, &pts[j].x, &pts[j].y)
       }
       // compute signed area
       var s float64
       for j := 0; j < 3; j++ {
           s += pts[j].Cross(pts[j+1])
       }
       s += pts[3].Cross(pts[0])
       area := s / 2
       if area < 0 {
           // reverse to make CCW
           for l, r := 0, 3; l < r; l, r = l+1, r-1 {
               pts[l], pts[r] = pts[r], pts[l]
           }
           area = -area
       }
       sumArea += area
       polygons[i] = pts
   }
   unionArea := polyUnion(polygons)
   fmt.Printf("%.9f\n", sumArea/unionArea)
}
