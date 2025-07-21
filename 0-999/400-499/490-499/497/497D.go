package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

type Point struct { x, y float64 }
type IPoint struct { x, y int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var Px, Py, Qx, Qy int64
   var n, m int
   // read P
   if _, err := fmt.Fscan(in, &Px, &Py); err != nil {
       return
   }
   fmt.Fscan(in, &n)
   A := make([]IPoint, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &A[i].x, &A[i].y)
   }
   // read Q
   fmt.Fscan(in, &Qx, &Qy)
   fmt.Fscan(in, &m)
   B := make([]IPoint, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &B[i].x, &B[i].y)
   }
   // vector D = Q - P
   Dx := Qx - Px
   Dy := Qy - Py
   // circle center C0 = -D
   Cx := float64(-Dx)
   Cy := float64(-Dy)
   R2i := Dx*Dx + Dy*Dy
   if R2i == 0 {
       fmt.Println("NO")
       return
   }
   R2 := float64(R2i)
   // vertex-vertex check
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           dx := A[i].x - B[j].x + Dx
           dy := A[i].y - B[j].y + Dy
           if dx*dx+dy*dy == R2i {
               fmt.Println("YES")
               return
           }
       }
   }
   // segment-circle intersection check for B vertices and A edges
   eps := 1e-9
   for j := 0; j < m; j++ {
       bvx := float64(B[j].x)
       bvy := float64(B[j].y)
       for i := 0; i < n; i++ {
           ni := (i + 1) % n
           // segment u = A[i] - B[j]
           u1x := float64(A[i].x) - bvx
           u1y := float64(A[i].y) - bvy
           u2x := float64(A[ni].x) - bvx
           u2y := float64(A[ni].y) - bvy
           dx := u2x - u1x
           dy := u2y - u1y
           a := dx*dx + dy*dy
           // relative to circle at C0
           ux := u1x - Cx
           uy := u1y - Cy
           if a < eps {
               // point
               if math.Abs(ux*ux+uy*uy-R2) < eps {
                   fmt.Println("YES")
                   return
               }
               continue
           }
           bq := 2 * (dx*ux + dy*uy)
           cq := ux*ux + uy*uy - R2
           disc := bq*bq - 4*a*cq
           if disc < 0 {
               continue
           }
           sd := math.Sqrt(disc)
           t1 := (-bq - sd) / (2 * a)
           t2 := (-bq + sd) / (2 * a)
           if (t1 >= -eps && t1 <= 1+eps) || (t2 >= -eps && t2 <= 1+eps) {
               fmt.Println("YES")
               return
           }
       }
   }
   // segment-circle for A vertices and B edges
   for i := 0; i < n; i++ {
       avx := float64(A[i].x)
       avy := float64(A[i].y)
       for j := 0; j < m; j++ {
           nj := (j + 1) % m
           u1x := avx - float64(B[j].x)
           u1y := avy - float64(B[j].y)
           u2x := avx - float64(B[nj].x)
           u2y := avy - float64(B[nj].y)
           dx := u2x - u1x
           dy := u2y - u1y
           a := dx*dx + dy*dy
           ux := u1x - Cx
           uy := u1y - Cy
           if a < eps {
               if math.Abs(ux*ux+uy*uy-R2) < eps {
                   fmt.Println("YES")
                   return
               }
               continue
           }
           bq := 2 * (dx*ux + dy*uy)
           cq := ux*ux + uy*uy - R2
           disc := bq*bq - 4*a*cq
           if disc < 0 {
               continue
           }
           sd := math.Sqrt(disc)
           t1 := (-bq - sd) / (2 * a)
           t2 := (-bq + sd) / (2 * a)
           if (t1 >= -eps && t1 <= 1+eps) || (t2 >= -eps && t2 <= 1+eps) {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}
