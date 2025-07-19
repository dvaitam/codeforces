package main

import (
   "fmt"
   "math"
)

const eps = 1e-8

// sign function
func sgn(x float64) int {
   if x > eps {
       return 1
   }
   if x < -eps {
       return -1
   }
   return 0
}

// Point
type Pt struct {
   x, y float64
}

func (p Pt) Len() float64 {
   return math.Hypot(p.x, p.y)
}

func Add(a, b Pt) Pt {
   return Pt{a.x + b.x, a.y + b.y}
}

func Sub(a, b Pt) Pt {
   return Pt{a.x - b.x, a.y - b.y}
}

func Dot(a, b Pt) float64 {
   return a.x*b.x + a.y*b.y
}

func Cross(a, b Pt) float64 {
   return a.x*b.y - a.y*b.x
}

// getInter checks intersection of segments p1-p2 and p3-p4, returns point and ok
func getInter(p1, p2, p3, p4 Pt) (bool, Pt) {
   v21 := Sub(p2, p1)
   v31 := Sub(p3, p1)
   v41 := Sub(p4, p1)
   d1 := Cross(v21, v31)
   d2 := Cross(v21, v41)
   v43 := Sub(p4, p3)
   v13 := Sub(p1, p3)
   v23 := Sub(p2, p3)
   d3 := Cross(v43, v13)
   d4 := Cross(v43, v23)
   s1 := sgn(d1)
   s2 := sgn(d2)
   s3 := sgn(d3)
   s4 := sgn(d4)
   if s2 == 0 {
       return false, Pt{}
   }
   var c Pt
   t := d2 - d1
   c.x = (p3.x*d2 - p4.x*d1) / t
   c.y = (p3.y*d2 - p4.y*d1) / t
   if s1*s2 <= 0 && s3*s4 <= 0 {
       return true, c
   }
   return false, Pt{}
}

// compute cost along polygon from id1+1 to id2 (mod n)
func getCost(pol []Pt, n, id1, id2 int, p1, p2 Pt) float64 {
   cost := Sub(pol[(id1+1)%n], p1).Len() + Sub(pol[id2], p2).Len()
   for i := (id1 + 1) % n; i != id2; i = (i + 1) % n {
       j := (i + 1) % n
       cost += Sub(pol[j], pol[i]).Len()
   }
   return cost
}

func getMinCost(src, dst Pt, pol []Pt, n int) float64 {
   // find intersections
   var vec [2]struct{
       id int
       p  Pt
   }
   cnt := 0
   // ensure pol[n] is same as pol[0]
   for i := 0; i < n; i++ {
       ok, c := getInter(src, dst, pol[i], pol[(i+1)%n])
       if ok {
           if cnt < 2 {
               vec[cnt] = struct{ id int; p Pt }{i, c}
           }
           cnt++
       }
   }
   if cnt != 2 {
       return Sub(dst, src).Len()
   }
   // order them
   if sgn(Dot(Sub(vec[0].p, src), Sub(vec[0].p, vec[1].p))) > 0 {
       vec[0], vec[1] = vec[1], vec[0]
   }
   // cost options
   interDist := Sub(vec[1].p, vec[0].p).Len()
   costInside := interDist * 2.0
   costCW := getCost(pol, n, vec[0].id, vec[1].id, vec[0].p, vec[1].p)
   costCCW := getCost(pol, n, vec[1].id, vec[0].id, vec[1].p, vec[0].p)
   base := Sub(vec[0].p, src).Len() + Sub(vec[1].p, dst).Len()
   return base + math.Min(costInside, math.Min(costCW, costCCW))
}

func main() {
   var src, dst Pt
   var n int
   if _, err := fmt.Scan(&src.x, &src.y, &dst.x, &dst.y); err != nil {
       return
   }
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   pol := make([]Pt, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&pol[i].x, &pol[i].y)
   }
   // compute
   res := getMinCost(src, dst, pol, n)
   fmt.Printf("%.8f\n", res)
}
