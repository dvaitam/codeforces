package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

const (
   ERROR    = 1e-9
   INFINITE = 1e9
)

func calTime(x, v []float64, i, j int) float64 {
   if math.Abs(v[i]-v[j]) > ERROR && math.Abs(x[j]-x[i]) > ERROR && (v[i]-v[j])*(x[j]-x[i]) > 0 {
       return (x[j] - x[i]) / (v[i] - v[j])
   }
   return INFINITE
}

func collid(v []float64, m []int, i, j int) {
   v1, v2 := v[i], v[j]
   m1, m2 := float64(m[i]), float64(m[j])
   v[i] = ((m1 - m2) * v1 + 2*m2*v2) / (m1 + m2)
   v[j] = ((m2 - m1) * v2 + 2*m1*v1) / (m1 + m2)
}

func passTime(x, v []float64, t float64) {
   for i := range x {
       x[i] += v[i] * t
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var T float64
   if _, err := fmt.Fscan(reader, &n, &T); err != nil {
       return
   }
   x := make([]float64, n)
   v := make([]float64, n)
   m := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &v[i], &m[i])
   }
   for T > 0 {
       colTime := INFINITE
       for i := 0; i < n; i++ {
           for j := i + 1; j < n; j++ {
               t := calTime(x, v, i, j)
               if t < colTime && t < T {
                   colTime = t
               }
           }
       }
       if colTime == INFINITE {
           passTime(x, v, T)
           break
       }
       passTime(x, v, colTime)
       // process all collisions at this time
       for i := 0; i < n; i++ {
           for j := i + 1; j < n; j++ {
               if math.Abs(x[i]-x[j]) < ERROR {
                   collid(v, m, i, j)
               }
           }
       }
       T -= colTime
   }
   for i := 0; i < n; i++ {
       fmt.Printf("%.6f\n", x[i])
   }
   fmt.Println()
}
