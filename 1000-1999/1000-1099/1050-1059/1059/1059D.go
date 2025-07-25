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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   xs := make([]float64, n)
   ys := make([]float64, n)
   var sign int
   for i := 0; i < n; i++ {
       var xi, yi float64
       fmt.Fscan(reader, &xi, &yi)
       xs[i] = xi
       ys[i] = yi
       s := 1
       if yi < 0 {
           s = -1
       }
       if i == 0 {
           sign = s
       } else if sign != s {
           fmt.Println(-1)
           return
       }
       if yi < 0 {
           ys[i] = -yi
       }
   }
   // All ys are now positive
   maxY := 0.0
   for _, y := range ys {
       if y > maxY {
           maxY = y
       }
   }
   // Minimum possible radius
   l := maxY / 2.0
   // feasibility check for given radius r
   feasible := func(r float64) bool {
       // x0 interval [L, R]
       L := -1e50
       R := 1e50
       for i := 0; i < n; i++ {
           y := ys[i]
           // delta = 2*y*r - y*y
           delta := 2.0*y*r - y*y
           if delta < 0 {
               return false
           }
           d := math.Sqrt(delta)
           xi := xs[i]
           lo := xi - d
           hi := xi + d
           if lo > L {
               L = lo
           }
           if hi < R {
               R = hi
           }
           if L > R {
               return false
           }
       }
       return L <= R
   }
   // Find an upper bound
   h := l
   if !feasible(h) {
       for {
           h *= 2.0
           if feasible(h) {
               break
           }
       }
   }
   // binary search
   for it := 0; it < 60; it++ {
       mid := (l + h) / 2.0
       if feasible(mid) {
           h = mid
       } else {
           l = mid
       }
   }
   fmt.Printf("%.10f\n", h)
}
