package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k, e, ix, iy int
   if _, err := fmt.Fscan(reader, &n, &k, &e, &ix, &iy); err != nil {
       return
   }
   ptsX := make([]int, n)
   ptsY := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &ptsX[i], &ptsY[i])
   }
   prob := make([]float64, n)
   prev := make([]float64, n+2)
   cur := make([]float64, n+2)
   // check function for given radius x
   chk := func(x float64) bool {
       for i := 0; i < n; i++ {
           dx := float64(ptsX[i] - ix)
           dy := float64(ptsY[i] - iy)
           d := math.Hypot(dx, dy)
           switch {
           case d < x:
               prob[i] = 1
           case d > x*1000:
               prob[i] = 0
           default:
               prob[i] = math.Exp(1 - (d*d)/(x*x))
           }
       }
       // initialize dp
       for i := range prev {
           prev[i] = 0
           cur[i] = 0
       }
       prev[0] = 1.0
       for i := 0; i < n; i++ {
           p := prob[i]
           // reset cur for this iteration
           for j := 0; j <= i+1; j++ {
               cur[j] = 0
           }
           for j := 0; j <= i; j++ {
               cur[j] += prev[j] * (1 - p)
               cur[j+1] += prev[j] * p
           }
           // swap prev and cur
           prev, cur = cur, prev
       }
       sum := 0.0
       limit := float64(e) / 1000.0
       for j := 0; j < k && j < len(prev); j++ {
           sum += prev[j]
       }
       return sum <= limit
   }
   // binary search for minimal x
   lo, hi := 0.0, 10000.0
   eps := 1e-7
   for hi-lo > eps {
       mid := (lo + hi) * 0.5
       if chk(mid) {
           hi = mid
       } else {
           lo = mid
       }
   }
   fmt.Printf("%.10f", lo)
}
