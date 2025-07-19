package main

import (
   "bufio"
   "fmt"
   "math"
   "math/rand"
   "os"
)

const eps = 1e-8

// sign returns -1 if x < -eps, 1 if x > eps, else 0
func sign(x float64) int {
   if x < -eps {
       return -1
   }
   if x > eps {
       return 1
   }
   return 0
}

// point3d represents a point in 3D space
type point3d struct { x, y, z float64 }

func (p point3d) add(q point3d) point3d {
   return point3d{p.x + q.x, p.y + q.y, p.z + q.z}
}

func (p point3d) sub(q point3d) point3d {
   return point3d{p.x - q.x, p.y - q.y, p.z - q.z}
}

func (p point3d) div(f float64) point3d {
   return point3d{p.x / f, p.y / f, p.z / f}
}

// len2 returns the squared length of the vector p
func (p point3d) len2() float64 {
   return p.x*p.x + p.y*p.y + p.z*p.z
}

// dis returns squared distance between points a and b
func dis(a, b point3d) float64 {
   d := a.sub(b)
   return d.len2()
}

// project projects point (x1,y1,z1) onto plane A*x+B*y+C*z=0
func project(A, B, C, x1, y1, z1 float64) point3d {
   lambda := (A*x1 + B*y1 + C*z1) / (A*A + B*B + C*C)
   return point3d{x1 - lambda*A, y1 - lambda*B, z1 - lambda*C}
}

// calc2 solves two linear equations:
// a*x + b*y + c = 0 and d*x + e*y + f = 0
// returns (x, y)
func calc2(a, b, c, d, e, f float64) (x, y float64) {
   if sign(a) == 0 {
       a, d = d, a
       b, e = e, b
       c, f = f, c
   }
   e += -b/a * d
   f += -c/a * d
   y = -f / e
   x = (-b*y - c) / a
   return
}

// cir_cen computes the center of the circle through points a, b, c
// lying in the plane with normal vector (A3,B3,C3)
func cir_cen(a, b, c point3d, A3, B3, C3 float64) point3d {
   var A1, B1, C1, D1 float64
   var A2, B2, C2, D2 float64
   var D3 float64
   A1 = -2*a.x + 2*b.x
   B1 = -2*a.y + 2*b.y
   C1 = -2*a.z + 2*b.z
   D1 = -(b.x*b.x + b.y*b.y + b.z*b.z - a.x*a.x - a.y*a.y - a.z*a.z)
   A2 = -2*c.x + 2*b.x
   B2 = -2*c.y + 2*b.y
   C2 = -2*c.z + 2*b.z
   D2 = -(b.x*b.x + b.y*b.y + b.z*b.z - c.x*c.x - c.y*c.y - c.z*c.z)
   D3 = 0
   if sign(A2) != 0 {
       A1, A2 = A2, A1
       B1, B2 = B2, B1
       C1, C2 = C2, C1
       D1, D2 = D2, D1
   }
   if sign(A3) != 0 {
       A1, A3 = A3, A1
       B1, B3 = B3, B1
       C1, C3 = C3, C1
       D1, D3 = D3, D1
   }
   if sign(A1) == 0 {
       return a
   }
   B2 += -B1/A1 * A2
   C2 += -C1/A1 * A2
   D2 += -D1/A1 * A2
   B3 += -B1/A1 * A3
   C3 += -C1/A1 * A3
   D3 += -D1/A1 * A3
   y, z := calc2(B2, C2, D2, B3, C3, D3)
   x := (-B1*y - C1*z - D1) / A1
   return point3d{x, y, z}
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   p := make([]point3d, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i].x, &p[i].y, &p[i].z)
   }
   for mm := 0; mm < m; mm++ {
       var A, B, C float64
       fmt.Fscan(in, &A, &B, &C)
       a := make([]point3d, n)
       for i := 0; i < n; i++ {
           a[i] = project(A, B, C, p[i].x, p[i].y, p[i].z)
       }
       rand.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
       o := point3d{0, 0, 0}
       var r2 float64
       for i := 0; i < n; i++ {
           if sign(dis(o, a[i]) - r2) > 0 {
               o = a[i]
               r2 = 0
               for j := 0; j < i; j++ {
                   if sign(dis(o, a[j]) - r2) > 0 {
                       o = a[i].add(a[j]).div(2)
                       r2 = o.sub(a[i]).len2()
                       for k := 0; k < j; k++ {
                           if sign(dis(o, a[k]) - r2) > 0 {
                               o = cir_cen(a[i], a[j], a[k], A, B, C)
                               r2 = o.sub(a[k]).len2()
                           }
                       }
                   }
               }
           }
       }
       fmt.Printf("%.9f\n", math.Sqrt(r2))
   }
}
