package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var cx, cy float64
   if _, err := fmt.Fscan(reader, &n, &cx, &cy); err != nil {
       return
   }
   xs := make([]float64, n)
   ys := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   best := 0.0
   worst := 1e25
   // compute max/min distance to vertices
   for i := 0; i < n; i++ {
       dx := xs[i] - cx
       dy := ys[i] - cy
       d := math.Hypot(dx, dy)
       if d > best {
           best = d
       }
       if d < worst {
           worst = d
       }
   }
   // compute min distance to edges (projection inside segment)
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       ax, ay := xs[i], ys[i]
       bx, by := xs[j], ys[j]
       // vectors
       abx := bx - ax
       aby := by - ay
       acx := cx - ax
       acy := cy - ay
       bcx := cx - bx
       bcy := cy - by
       // dot products
       if abx*acx+aby*acy >= 0 && (-abx)*bcx+(-aby)*bcy >= 0 {
           // cross product magnitude
           cross := math.Abs(abx*bcy - aby*bcx)
           // edge length
           elen := math.Hypot(abx, aby)
           if elen > 0 {
               d := cross / elen
               if d < worst {
                   worst = d
               }
           }
       }
   }
   // area difference: pi*(best^2 - worst^2)
   area := math.Pi * (best*best - worst*worst)
   fmt.Printf("%.12f\n", area)
}
