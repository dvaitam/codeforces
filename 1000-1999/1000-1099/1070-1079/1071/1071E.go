package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var w, h float64
   if _, err := fmt.Fscan(in, &n, &w, &h); err != nil {
       return
   }
   var e1, e2 float64
   fmt.Fscan(in, &e1, &e2)
   t := make([]float64, n)
   x := make([]float64, n)
   y := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &t[i], &x[i], &y[i])
   }
   // Precompute allowed slope intervals M_low, M_high
   Mlow := make([]float64, n)
   Mhigh := make([]float64, n)
   for i := 0; i < n; i++ {
       yi := y[i]
       yim := h - yi
       // from p bounds: (x-w)/y <= m <= x/y
       lo1 := (x[i] - w) / yi
       hi1 := x[i] / yi
       // from q bounds: -x/(h-y) <= m <= (w-x)/(h-y)
       lo2 := -x[i] / yim
       hi2 := (w - x[i]) / yim
       lo := math.Max(lo1, lo2)
       hi := math.Min(hi1, hi2)
       Mlow[i] = lo
       Mhigh[i] = hi
       if lo > hi {
           fmt.Println(-1)
           return
       }
   }
   // initial slope
   m0 := (e2 - e1) / h
   // check feasibility for large v
   const INF = 1e8
   if !check(INF, n, w, h, e1, e2, t, x, y, Mlow, Mhigh, m0) {
       fmt.Println(-1)
       return
   }
   // binary search v
   lo := 0.0
   hi := INF
   for it := 0; it < 60; it++ {
       mid := (lo + hi) / 2
       if check(mid, n, w, h, e1, e2, t, x, y, Mlow, Mhigh, m0) {
           hi = mid
       } else {
           lo = mid
       }
   }
   fmt.Printf("%.10f\n", hi)
}

// check if speed v is enough
func check(v float64, n int, w, h, e1, e2 float64,
   t, x, y, Mlow, Mhigh []float64, m0 float64) bool {
   // first drop i=0
   dt := t[0]
   D := v * dt
   // f constraint: (x0 - e1 - D)/y0 <= m <= (x0 - e1 + D)/y0
   f1 := (x[0] - e1 - D) / y[0]
   f2 := (x[0] - e1 + D) / y[0]
   // g constraint: (e2 - x0 - D)/(h-y0) <= m <= (e2 - x0 + D)/(h-y0)
   g1 := (e2 - x[0] - D) / (h - y[0])
   g2 := (e2 - x[0] + D) / (h - y[0])
   L := math.Max(Mlow[0], math.Max(f1, g1))
   R := math.Min(Mhigh[0], math.Min(f2, g2))
   if L > R {
       return false
   }
   // iterate over subsequent drops
   for i := 1; i < n; i++ {
       prevx := x[i-1]
       prevy := y[i-1]
       dt = t[i] - t[i-1]
       D = v * dt
       yi := y[i]
       yim := h - yi
       // f interval union
       a := (x[i] - prevx - D) / yi
       b := (x[i] - prevx + D) / yi
       f_low := a + L*(prevy/yi)
       f_high := b + R*(prevy/yi)
       // g interval union
       c := (prevx - x[i] - D) / yim
       d := (prevx - x[i] + D) / yim
       g_low := c + L*((h-prevy)/yim)
       g_high := d + R*((h-prevy)/yim)
       // intersect with Mlow/high
       L = math.Max(Mlow[i], math.Max(f_low, g_low))
       R = math.Min(Mhigh[i], math.Min(f_high, g_high))
       if L > R {
           return false
       }
   }
   return true
}
