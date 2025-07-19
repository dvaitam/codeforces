package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var v int
   if _, err := fmt.Fscan(in, &n, &v); err != nil {
       return
   }
   al := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &al[i])
   }
   var m int
   fmt.Fscan(in, &m)
   type obstacle struct{ dist, height float64 }
   W := make([]obstacle, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &W[i].dist, &W[i].height)
   }
   sort.Slice(W, func(i, j int) bool {
       return W[i].dist < W[j].dist
   })
   d := make([]int, n)
   for i := 0; i < n; i++ {
       d[i] = i
   }
   sort.Slice(d, func(i, j int) bool {
       return al[d[i]] < al[d[j]]
   })
   X := make([]float64, n)
   Y := make([]float64, n)
   j := 0
   vf := float64(v)
   const g = 9.8
   for _, idx := range d {
       alpha := al[idx]
       maxRange := vf*vf/g*math.Sin(2*alpha)
       hit := false
       for j < m && W[j].dist < maxRange {
           t := W[j].dist / (vf * math.Cos(alpha))
           h := vf*math.Sin(alpha)*t - g*t*t/2
           if h <= W[j].height {
               X[idx] = W[j].dist
               Y[idx] = h
               hit = true
               break
           }
           j++
       }
       if !hit {
           X[idx] = maxRange
           Y[idx] = 0
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < n; i++ {
       fmt.Fprintf(out, "%.6f %.6f\n", X[i], Y[i])
   }
}
