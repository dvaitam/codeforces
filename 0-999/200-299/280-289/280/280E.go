package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func clamp(x, lo, hi float64) float64 {
   if x < lo {
       return lo
   }
   if x > hi {
       return hi
   }
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var q, a, b float64
   if _, err := fmt.Fscan(in, &n, &q, &a, &b); err != nil {
       return
   }
   xs := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &xs[i])
   }
   y := make([]float64, n)
   L := make([]float64, n)
   U := make([]float64, n)
   for i := 0; i < n; i++ {
       L[i] = 1 + float64(i)*a
       U[i] = q - float64(n-1-i)*a
   }
   // init
   for i := 0; i < n; i++ {
       y[i] = clamp(xs[i], L[i], U[i])
   }
   // iterative projections
   for it := 0; it < 200; it++ {
       maxd := 0.0
       // forward diff constraints
       for i := 0; i+1 < n; i++ {
           lo := y[i] + a
           hi := y[i] + b
           old := y[i+1]
           y[i+1] = clamp(y[i+1], lo, hi)
           if d := math.Abs(y[i+1] - old); d > maxd {
               maxd = d
           }
       }
       // backward diff constraints
       for i := n - 1; i > 0; i-- {
           hi := y[i] - a
           lo := y[i] - b
           old := y[i-1]
           y[i-1] = clamp(y[i-1], math.Max(lo, L[i-1]), math.Min(hi, U[i-1]))
           if d := math.Abs(y[i-1] - old); d > maxd {
               maxd = d
           }
       }
       if maxd < 1e-9 {
           break
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < n; i++ {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprintf(out, "%.9f", y[i])
   }
   out.WriteByte('\n')
   cost := 0.0
   for i := 0; i < n; i++ {
       d := y[i] - xs[i]
       cost += d * d
   }
   fmt.Fprintf(out, "%.9f", cost)
}
