package main

import (
   "fmt"
   "math"
)

const eps = 1e-8

// Point
type P struct {
   x, y float64
}

func (p P) Sub(a P) P  { return P{p.x - a.x, p.y - a.y} }
func (p P) Add(a P) P  { return P{p.x + a.x, p.y + a.y} }
func (p P) Mul(k float64) P { return P{p.x * k, p.y * k} }
func (p P) Len() float64  { return math.Hypot(p.x, p.y) }

// Circle
type C struct {
   x, y, r float64
}

func (c C) O() P { return P{c.x, c.y} }

// Line: ax + by + c = 0
type L struct {
   a, b, c float64
}

// CL: wrapper for circle or line
type CL struct {
   typ int // 0=circle, 1=line
   c   C
   l   L
}

// cross of two lines
func crossLL(l1, l2 L) []P {
   det := l1.a*l2.b - l1.b*l2.a
   if math.Abs(det) < eps {
       return nil
   }
   det1 := -(l1.c*l2.b - l1.b*l2.c)
   det2 := -(l1.a*l2.c - l1.c*l2.a)
   return []P{{det1 / det, det2 / det}}
}

// cross circle and line
func crossCL(c C, l L) []P {
   var res []P
   al, be := l.b, -l.a
   var x0, y0 float64
   if math.Abs(l.a) < math.Abs(l.b) {
       x0 = 0
       y0 = -l.c / l.b
   } else {
       y0 = 0
       x0 = -l.c / l.a
   }
   A := al*al + be*be
   B := 2*al*(x0-c.x) + 2*be*(y0-c.y)
   Cq := (x0-c.x)*(x0-c.x) + (y0-c.y)*(y0-c.y) - c.r*c.r
   D := B*B - 4*A*Cq
   if D < -eps {
       return nil
   }
   if D < 0 {
       D = 0
   }
   t1 := (-B + math.Sqrt(D)) / (2 * A)
   res = append(res, P{x0 + al*t1, y0 + be*t1})
   t2 := (-B - math.Sqrt(D)) / (2 * A)
   res = append(res, P{x0 + al*t2, y0 + be*t2})
   return res
}

// cross generic CL
func crossCLgen(cl1, cl2 CL) []P {
   if cl1.typ == 0 && cl2.typ == 0 {
       // circle-circle
       c1, c2 := cl1.c, cl2.c
       a := 2*(c2.x - c1.x)
       b := 2*(c2.y - c1.y)
       c0 := c2.r*c2.r - c1.r*c1.r + c1.x*c1.x - c2.x*c2.x + c1.y*c1.y - c2.y*c2.y
       return crossCL(c1, L{a, b, c0})
   }
   if cl1.typ == 0 && cl2.typ == 1 {
       return crossCL(cl1.c, cl2.l)
   }
   if cl1.typ == 1 && cl2.typ == 0 {
       return crossCL(cl2.c, cl1.l)
   }
   if cl1.typ == 1 && cl2.typ == 1 {
       return crossLL(cl1.l, cl2.l)
   }
   return nil
}

// get line for radical axis of two circles
func getL(c1, c2 C) CL {
   a := 2*c2.x - 2*c1.x
   b := 2*c2.y - 2*c1.y
   c0 := c1.x*c1.x - c2.x*c2.x + c1.y*c1.y - c2.y*c2.y
   return CL{typ: 1, l: L{a, b, c0}}
}

// get circle for coaxal circles
func getC(c1, c2 C) CL {
   if c1.r > c2.r {
       return getC(c2, c1)
   }
   cr := c1.r / c2.r
   o1 := c1.O()
   o2 := c2.O()
   v := o2.Sub(o1)
   p1 := o1.Add(v.Mul(cr / (1 + cr)))
   p2 := o1.Add(v.Mul(cr / (cr - 1)))
   o := p1.Add(p2).Mul(0.5)
   r := p1.Sub(o).Len()
   return CL{typ: 0, c: C{o.x, o.y, r}}
}

// choose line or circle
func getCL(c1, c2 C) CL {
   if math.Abs(c1.r-c2.r) < eps {
       return getL(c1, c2)
   }
   return getC(c1, c2)
}

func main() {
   var c [3]C
   for i := 0; i < 3; i++ {
       if _, err := fmt.Scan(&c[i].x, &c[i].y, &c[i].r); err != nil {
           return
       }
   }
   cl1 := getCL(c[0], c[1])
   cl2 := getCL(c[1], c[2])
   cr := crossCLgen(cl1, cl2)
   mi := 1e100
   var ans P
   for _, p := range cr {
       var q [3]float64
       ok := true
       for j := 0; j < 3; j++ {
           q[j] = p.Sub(c[j].O()).Len() / c[j].r
           if math.Abs(q[j]-q[0]) > eps {
               ok = false
           }
       }
       if q[0] < 1-eps {
           ok = false
       }
       if !ok {
           continue
       }
       if q[0] < mi {
           mi = q[0]
           ans = p
       }
   }
   if mi < 1e50 {
       fmt.Printf("%.5f %.5f\n", ans.x, ans.y)
   }
}
