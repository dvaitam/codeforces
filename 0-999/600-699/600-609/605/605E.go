package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([][]float64, n)
   for i := 0; i < n; i++ {
       a[i] = make([]float64, n)
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           a[i][j] = float64(x) / 100.0
       }
   }

   v := make([]float64, n)
   cnt := make([]float64, n)
   f := make([]float64, n)
   used := make([]bool, n)
   inf := math.Inf(1)
   for i := 0; i < n; i++ {
       f[i] = inf
       cnt[i] = 1.0
   }
   f[n-1] = 0.0

   for t := 0; t < n; t++ {
       k := -1
       for i := 0; i < n; i++ {
           if used[i] {
               continue
           }
           if k == -1 || f[i] < f[k] {
               k = i
           }
       }
       used[k] = true
       for i := 0; i < n; i++ {
           if used[i] {
               continue
           }
           v[i] += cnt[i] * a[i][k] * (f[k] + 1.0)
           cnt[i] *= (1.0 - a[i][k])
           if cnt[i] != 1.0 {
               val := (v[i] + cnt[i]) / (1.0 - cnt[i])
               if val < f[i] {
                   f[i] = val
               }
           }
       }
   }

   fmt.Fprintf(writer, "%.12f\n", f[0])
}
