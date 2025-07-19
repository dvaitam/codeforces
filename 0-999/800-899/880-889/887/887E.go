package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var x1, y1, x2, y2 float64
   if _, err := fmt.Fscan(in, &x1, &y1, &x2, &y2); err != nil {
       return
   }
   // midpoint
   mx := (x1 + x2) / 2
   my := (y1 + y2) / 2
   // vector from p2 to p1
   dx := x1 - x2
   dy := y1 - y2
   // half distance
   w := math.Hypot(dx, dy) / 2
   // perpendicular unit vector v = (-dy, dx)
   vx := -dy
   vy := dx
   norm := math.Hypot(vx, vy)
   vx /= norm
   vy /= norm

   var n int
   fmt.Fscan(in, &n)
   type pair struct{ x float64; y int }
   q := make([]pair, 0, 2*n+1)
   // initial dummy
   q = append(q, pair{0, 0})

   // get projection coordinate for given point and threshold val
   get := func(px, py, val float64) float64 {
       // function d(len) = distance(p, mid+v*len) - sqrt(w^2 + len^2)
       d := func(t float64) float64 {
           cx := mx + vx*t
           cy := my + vy*t
           dist := math.Hypot(px-cx, py-cy)
           return dist - math.Sqrt(w*w+t*t)
       }
       l, r := -1e12, 1e12
       fl := d(l) > val
       for i := 0; i < 100; i++ {
           mid := (l + r) / 2
           f := d(mid) > val
           if f == fl {
               l = mid
           } else {
               r = mid
           }
       }
       return r
   }

   for i := 0; i < n; i++ {
       var px, py, R float64
       fmt.Fscan(in, &px, &py, &R)
       l := get(px, py, R)
       rr := get(px, py, -R)
       if l > rr {
           l, rr = rr, l
       }
       q = append(q, pair{l, 1})
       q = append(q, pair{rr, -1})
   }
   sort.Slice(q, func(i, j int) bool {
       if q[i].x == q[j].x {
           return q[i].y < q[j].y
       }
       return q[i].x < q[j].x
   })
   s := 0
   mn := 1e12
   for _, p := range q {
       if s == 0 {
           mn = math.Min(mn, math.Abs(p.x))
       }
       s += p.y
       if s == 0 {
           mn = math.Min(mn, math.Abs(p.x))
       }
   }
   ans := math.Sqrt(mn*mn + w*w)
   // output with high precision
   fmt.Printf("%.20f\n", ans)
}
