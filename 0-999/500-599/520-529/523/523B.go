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

   var n, t, m int
   var c float64
   if _, err := fmt.Fscan(reader, &n, &t, &c); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   realAvg := make([]float64, n+1)
   approx := make([]float64, n+1)
   var sum int64
   ft := float64(t)
   for i := 1; i <= n; i++ {
       sum += int64(a[i])
       if i >= t {
           sum -= int64(a[i-t])
           realAvg[i] = float64(sum) / ft
       }
       approx[i] = (approx[i-1] + float64(a[i])/ft) / c
   }
   fmt.Fscan(reader, &m)
   for j := 0; j < m; j++ {
       var x int
       fmt.Fscan(reader, &x)
       r := realAvg[x]
       ap := approx[x]
       errRel := math.Abs(ap-r) / r
       fmt.Fprintf(writer, "%.5f %.5f %.5f", r, ap, errRel)
       if j < m-1 {
           writer.WriteByte('\n')
       }
   }
}
