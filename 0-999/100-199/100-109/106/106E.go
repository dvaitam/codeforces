package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const eps = 1e-12

// point3D represents a point in 3D space
type point3D struct { x, y, z float64 }

var (
   p      []point3D       // input points
   outer  [4]point3D      // boundary points for sphere
   ret    point3D         // current sphere center
   radius float64         // current sphere squared radius
   nouter int             // number of boundary points
)

// dissqr3D returns squared distance between two points
func dissqr3D(a, b *point3D) float64 {
   dx := a.x - b.x
   dy := a.y - b.y
   dz := a.z - b.z
   return dx*dx + dy*dy + dz*dz
}

// dot returns dot product of two vectors
func dot(a, b *point3D) float64 {
   return a.x*b.x + a.y*b.y + a.z*b.z
}

// det2 computes determinant of 2x2 matrix
func det2(m [2][2]float64) float64 {
   return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}

// det3 computes determinant of 3x3 matrix
func det3(m [3][3]float64) float64 {
   return m[0][0]*m[1][1]*m[2][2] + m[0][1]*m[1][2]*m[2][0] + m[0][2]*m[2][1]*m[1][0] -
          m[0][2]*m[1][1]*m[2][0] - m[0][1]*m[1][0]*m[2][2] - m[0][0]*m[1][2]*m[2][1]
}

// ball computes the minimal sphere for current outer points
func ball() {
   // initialize
   ret = point3D{0, 0, 0}
   radius = 0
   switch nouter {
   case 1:
       ret = outer[0]
   case 2:
       ret.x = (outer[0].x + outer[1].x) / 2
       ret.y = (outer[0].y + outer[1].y) / 2
       ret.z = (outer[0].z + outer[1].z) / 2
       radius = dissqr3D(&ret, &outer[0])
   case 3:
       var q [2]point3D
       for i := 0; i < 2; i++ {
           q[i].x = outer[i+1].x - outer[0].x
           q[i].y = outer[i+1].y - outer[0].y
           q[i].z = outer[i+1].z - outer[0].z
       }
       var m [2][2]float64
       var sol [2]float64
       for i := 0; i < 2; i++ {
           for j := 0; j < 2; j++ {
               m[i][j] = dot(&q[i], &q[j]) * 2
           }
           sol[i] = dot(&q[i], &q[i])
       }
       d := det2(m)
       if math.Abs(d) < eps {
           return
       }
       var L [2]float64
       L[0] = (sol[0]*m[1][1] - sol[1]*m[0][1]) / d
       L[1] = (sol[1]*m[0][0] - sol[0]*m[1][0]) / d
       ret.x = outer[0].x + q[0].x*L[0] + q[1].x*L[1]
       ret.y = outer[0].y + q[0].y*L[0] + q[1].y*L[1]
       ret.z = outer[0].z + q[0].z*L[0] + q[1].z*L[1]
       radius = dissqr3D(&ret, &outer[0])
   case 4:
       var q [3]point3D
       var sol [3]float64
       var m [3][3]float64
       for i := 0; i < 3; i++ {
           q[i].x = outer[i+1].x - outer[0].x
           q[i].y = outer[i+1].y - outer[0].y
           q[i].z = outer[i+1].z - outer[0].z
           sol[i] = dot(&q[i], &q[i])
       }
       for i := 0; i < 3; i++ {
           for j := 0; j < 3; j++ {
               m[i][j] = dot(&q[i], &q[j]) * 2
           }
       }
       d := det3(m)
       if math.Abs(d) < eps {
           return
       }
       var L [3]float64
       for j := 0; j < 3; j++ {
           // replace column j with sol
           for i := 0; i < 3; i++ {
               m[i][j] = sol[i]
           }
           L[j] = det3(m) / d
           // restore column j
           for i := 0; i < 3; i++ {
               m[i][j] = dot(&q[i], &q[j]) * 2
           }
       }
       ret = outer[0]
       for i := 0; i < 3; i++ {
           ret.x += q[i].x * L[i]
           ret.y += q[i].y * L[i]
           ret.z += q[i].z * L[i]
       }
       radius = dissqr3D(&ret, &outer[0])
   }
}

// minball implements Welzl's algorithm recursively
func minball(n int) {
   ball()
   if nouter < 4 {
       for i := 0; i < n; i++ {
           if dissqr3D(&ret, &p[i]) - radius > eps {
               outer[nouter] = p[i]
               nouter++
               minball(i)
               nouter--
               if i > 0 {
                   tmp := p[i]
                   copy(p[1:i+1], p[0:i])
                   p[0] = tmp
               }
           }
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           return
       }
       p = make([]point3D, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &p[i].x, &p[i].y, &p[i].z)
       }
       nouter = 0
       minball(n)
       fmt.Printf("%.8f %.8f %.8f\n", ret.x+eps, ret.y+eps, ret.z+eps)
   }
}
