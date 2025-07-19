package main

import (
   "fmt"
   "math"
)

// Point represents a 2D point or vector
type Point struct {
   x, y float64
}

// add returns the vector sum of p and q
func (p Point) add(q Point) Point {
   return Point{p.x + q.x, p.y + q.y}
}

// mul returns the scalar multiplication of p by s
func (p Point) mul(s float64) Point {
   return Point{p.x * s, p.y * s}
}

// dist returns the Euclidean distance between p and q
func dist(p, q Point) float64 {
   return math.Hypot(p.x-q.x, p.y-q.y)
}

func main() {
   var t1, t2 float64
   if _, err := fmt.Scan(&t1, &t2); err != nil {
       return
   }
   var a, b, c Point
   // note: input order is a, c, b as in original code
   if _, err := fmt.Scan(&a.x, &a.y, &c.x, &c.y, &b.x, &b.y); err != nil {
       return
   }
   ab := dist(a, b)
   bc := dist(b, c)
   ac := dist(a, c)
   t1 += ab + bc + 1e-12
   t2 += ac + 1e-12
   // simple case
   if ab+bc < t2 {
       fmt.Printf("%.10f\n", math.Min(t1, t2))
       return
   }

   // cal computes the minimal feasible time at parameter k
   cal := func(k float64) float64 {
       // p moves from b to c
       p := b.mul(1-k).add(c.mul(k))
       ap := dist(a, p)
       // check if no intermediate stop needed
       if ap+(k+1)*bc < t1 && ap+(1-k)*bc < t2 {
           return math.Min(t1-(k+1)*bc, t2-(1-k)*bc)
       }
       // binary search on parameter m between a and p
       l, r := 0.0, 1.0
       for r-l > 1e-15 {
           m := (l + r) / 2
           p1 := a.mul(1-m).add(p.mul(m))
           if ap*m+dist(p1, b)+bc < t1 && ap*m+dist(p1, c) < t2 {
               l = m
           } else {
               r = m
           }
       }
       return ((l + r) / 2) * ap
   }

   // ternary search on k in [0,1]
   l, r := 0.0, 1.0
   for r-l > 1e-15 {
       m1 := (2*l + r) / 3
       m2 := (l + 2*r) / 3
       if cal(m1)-cal(m2) < 1e-12 {
           l = m1
       } else {
           r = m2
       }
   }
   ans := cal((l + r) / 2)
   fmt.Printf("%.10f\n", ans)
}
