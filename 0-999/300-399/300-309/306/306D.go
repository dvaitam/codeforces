package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// Vec represents a 2D vector or complex number
type Vec struct {
   x, y float64
}

// add returns the vector sum of v and u
func (v Vec) add(u Vec) Vec {
   return Vec{v.x + u.x, v.y + u.y}
}

// mul returns the complex multiplication of v and u
func (v Vec) mul(u Vec) Vec {
   return Vec{v.x*u.x - v.y*u.y, v.x*u.y + v.y*u.x}
}

// abs returns the magnitude of v
func (v Vec) abs() float64 {
   return math.Hypot(v.x, v.y)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n < 5 {
       fmt.Fprintln(writer, "No solution")
       return
   }

   const base = 690.0
   PI := math.Acos(-1.0)
   // rotation by 2*PI/n
   u := Vec{math.Cos(2*PI/float64(n)), math.Sin(2*PI/float64(n))}
   v := Vec{base, 0}
   // 1-indexed points
   pts := make([]Vec, n+2)
   pts[1] = Vec{0, 0}

   // build first n-1 points
   for i := 2; i <= n-1; i++ {
       pts[i] = pts[i-1].add(v)
       // rotate v
       v = v.mul(u)
       // scale to keep length reduced by 0.01
       r := v.abs()
       if r != 0 {
           factor := (r - 0.01) / r
           v.x *= factor
           v.y *= factor
       }
   }

   // solve for last point
   v1 := v
   v2 := v.mul(u)
   a1, b1 := v1.x, v2.x
   c1 := pts[1].x - pts[n-1].x
   a2, b2 := v1.y, v2.y
   c2 := pts[1].y - pts[n-1].y
   det := a1*b2 - a2*b1
   xx1 := (c1*b2 - c2*b1) / det
   // yy1 := (a1*c2 - a2*c1) / det  // unused
   // last point
   pts[n] = pts[n-1].add(Vec{v1.x * xx1, v1.y * xx1})

   // output
   for i := 1; i <= n; i++ {
       fmt.Fprintf(writer, "%.9f %.9f\n", pts[i].x, pts[i].y)
   }
}
