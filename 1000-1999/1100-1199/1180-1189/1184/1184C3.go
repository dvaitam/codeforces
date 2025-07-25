package main

import (
   "bufio"
   "fmt"
   "math"
   "math/rand"
   "os"
   "time"
)

type Point struct { x, y float64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var k, n int
   fmt.Fscan(in, &k)
   fmt.Fscan(in, &n)
   pts := make([]Point, k*n)
   for i := range pts {
       var xi, yi int
       fmt.Fscan(in, &xi, &yi)
       pts[i] = Point{float64(xi), float64(yi)}
   }
   rand.Seed(time.Now().UnixNano())
   circles := recoverRings(k, n, pts)
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for _, c := range circles {
       fmt.Fprintf(out, "%.10f %.10f %.10f\n", c.x, c.y, c.r)
   }
}

type Circle struct { x, y, r float64 }

// recoverRings finds k circles from points
func recoverRings(k, n int, pts []Point) []Circle {
   remain := make([]Point, len(pts))
   copy(remain, pts)
   var res []Circle
   // threshold for inliers
   thrCount := int(0.6 * float64(n))
   for ri := 0; ri < k; ri++ {
       bestCnt := 0
       var bestC Circle
       var bestInliers []int
       // RANSAC
       maxIter := 2000
       N := len(remain)
       for it := 0; it < maxIter; it++ {
           // pick 3 distinct
           i := rand.Intn(N)
           j := rand.Intn(N-1); if j >= i { j++ }
           l := rand.Intn(N-2)
           if l >= min(i, j) {
               l++
               if l >= max(i, j) { l++ }
           }
           p1, p2, p3 := remain[i], remain[j], remain[l]
           cx, cy, cr, ok := fitCircle(p1, p2, p3)
           if !ok || cr <= 0 { continue }
           tol := 0.1*cr + 1.0
           var inliers []int
           cnt := 0
           for idx, p := range remain {
               d := math.Hypot(p.x-cx, p.y-cy)
               if math.Abs(d-cr) <= tol {
                   cnt++
                   inliers = append(inliers, idx)
               }
           }
           if cnt > bestCnt {
               bestCnt = cnt
               bestC = Circle{cx, cy, cr}
               bestInliers = inliers
               if bestCnt >= thrCount {
                   break
               }
           }
       }
       // refine with Kasa fit on best inliers
       if len(bestInliers) > 3 {
           sub := make([]Point, len(bestInliers))
           for i, idx := range bestInliers {
               sub[i] = remain[idx]
           }
           cx, cy, cr := fitCircleKasa(sub)
           bestC = Circle{cx, cy, cr}
       }
       res = append(res, bestC)
       // remove inliers
       mark := make([]bool, len(remain))
       for _, idx := range bestInliers {
           mark[idx] = true
       }
       newRem := remain[:0]
       for i, p := range remain {
           if !mark[i] {
               newRem = append(newRem, p)
           }
       }
       remain = newRem
   }
   return res
}

// fitCircle returns center and radius from 3 points; ok=false if degenerate
func fitCircle(p1, p2, p3 Point) (cx, cy, r float64, ok bool) {
   x1, y1 := p1.x, p1.y
   x2, y2 := p2.x, p2.y
   x3, y3 := p3.x, p3.y
   d := 2 * (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2))
   if math.Abs(d) < 1e-8 {
       return 0, 0, 0, false
   }
   sq1 := x1*x1 + y1*y1
   sq2 := x2*x2 + y2*y2
   sq3 := x3*x3 + y3*y3
   ux := (sq1*(y2-y3) + sq2*(y3-y1) + sq3*(y1-y2)) / d
   uy := (sq1*(x3-x2) + sq2*(x1-x3) + sq3*(x2-x1)) / d
   r = math.Hypot(x1-ux, y1-uy)
   return ux, uy, r, true
}

// fitCircleKasa fits circle by least-squares (Kasa method)
func fitCircleKasa(pts []Point) (cx, cy, r float64) {
   m := float64(len(pts))
   var mx, my float64
   for _, p := range pts {
       mx += p.x
       my += p.y
   }
   mx /= m; my /= m
   var Suu, Suv, Svv, Suz, Svz float64
   for _, p := range pts {
       X := p.x - mx
       Y := p.y - my
       Z := X*X + Y*Y
       Suu += X * X
       Suv += X * Y
       Svv += Y * Y
       Suz += X * Z
       Svz += Y * Z
   }
   det := Suu*Svv - Suv*Suv
   var a, b float64
   if math.Abs(det) > 1e-12 {
       a = (Svv*Suz - Suv*Svz) / det
       b = (Suu*Svz - Suv*Suz) / det
   }
   cx = mx + a/2
   cy = my + b/2
   // average radius
   var sumr float64
   for _, p := range pts {
       sumr += math.Hypot(p.x-cx, p.y-cy)
   }
   r = sumr / m
   return
}
