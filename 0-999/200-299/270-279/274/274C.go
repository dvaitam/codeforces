package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const ep = 1e-9

// point
type pnt struct {
   x, y float64
}

func (a pnt) sub(b pnt) pnt   { return pnt{a.x - b.x, a.y - b.y} }
func (a pnt) add(b pnt) pnt   { return pnt{a.x + b.x, a.y + b.y} }
func (a pnt) mul(s float64) pnt { return pnt{a.x * s, a.y * s} }
func (a pnt) div(s float64) pnt { return pnt{a.x / s, a.y / s} }
func (a pnt) dot(b pnt) float64 { return a.x*b.x + a.y*b.y }
func (a pnt) cross(b pnt) float64 { return a.x*b.y - a.y*b.x }
func (a pnt) dist() float64 { return math.Hypot(a.x, a.y) }

func outercenter(a, b, c pnt) pnt {
   c1 := (a.dot(a) - b.dot(b)) / 2
   c2 := (a.dot(a) - c.dot(c)) / 2
   // denom = (a-b) cross (a-c)
   d1 := a.sub(b).cross(a.sub(c))
   // another denom variant
   // compute x0, y0
   x0 := (c1*(a.y-c.y) - c2*(a.y-b.y)) / d1
   y0 := (c1*(a.x-c.x) - c2*(a.x-b.x)) / -d1
   return pnt{x0, y0}
}

func chk(a, b, c pnt, pts []pnt, i, j, k int) bool {
   // circle with diameter bc, center o, radius r
   r := b.sub(c).dist() / 2
   o := b.add(c).div(2)
   bf := false
   for t := range pts {
       if t == i || t == j || t == k {
           continue
       }
       // check if point on opposite side and within twice radius
       if (pts[t].sub(c).cross(b.sub(c)) * (a.sub(c).cross(b.sub(c)))) < -ep {
           if pts[t].sub(b).dist() < 2*r && pts[t].sub(c).dist() < 2*r {
               bf = true
           }
       }
       if o.sub(pts[t]).dist() < r-ep {
           return false
       }
   }
   return bf
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   pts := make([]pnt, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
   }
   ans := -1.0
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           for k := j + 1; k < n; k++ {
               fijk := pts[i].sub(pts[j]).dot(pts[k].sub(pts[j]))
               fikj := pts[i].sub(pts[k]).dot(pts[j].sub(pts[k]))
               fjik := pts[j].sub(pts[i]).dot(pts[k].sub(pts[i]))
               if fijk < -ep || fikj < -ep || fjik < -ep {
                   continue
               }
               // colinear cases
               if math.Abs(fijk) < ep {
                   // j-i-k colinear, diameter bc where b=pts[i], c=pts[k]? mimic C++ logic
                   a, b, c := pts[j], pts[i], pts[k]
                   if chk(a, b, c, pts, i, j, k) {
                       r := b.sub(c).dist() / 2
                       if r > ans {
                           ans = r
                       }
                   }
               } else if math.Abs(fikj) < ep {
                   a, b, c := pts[k], pts[i], pts[j]
                   if chk(a, b, c, pts, i, j, k) {
                       r := b.sub(c).dist() / 2
                       if r > ans {
                           ans = r
                       }
                   }
               } else if math.Abs(fjik) < ep {
                   a, b, c := pts[i], pts[j], pts[k]
                   if chk(a, b, c, pts, i, j, k) {
                       r := b.sub(c).dist() / 2
                       if r > ans {
                           ans = r
                       }
                   }
               } else {
                   // circumcircle
                   area := math.Abs(pts[i].sub(pts[j]).cross(pts[k].sub(pts[j])))
                   r := pts[i].sub(pts[j]).dist() * pts[i].sub(pts[k]).dist() * pts[j].sub(pts[k]).dist() / area / 2
                   o := outercenter(pts[i], pts[j], pts[k])
                   ok := true
                   for t := range pts {
                       if t == i || t == j || t == k {
                           continue
                       }
                       if o.sub(pts[t]).dist() < r-ep {
                           ok = false
                           break
                       }
                   }
                   if ok && r > ans {
                       ans = r
                   }
               }
           }
       }
   }
   if ans < -ep {
       fmt.Println(-1)
   } else {
       fmt.Printf("%.12f\n", ans)
   }
}
