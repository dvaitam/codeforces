package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func max(a, b float64) float64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   w := make([]int, n)
   v := make([]float64, n)
   sumW := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i], &v[i])
       sumW += w[i]
   }
   // dp[j][k]: max sum v for choosing j glasses with total width k
   const INF = 1e18
   dp := make([][]float64, n+1)
   for j := 0; j <= n; j++ {
       dp[j] = make([]float64, sumW+1)
       for k := 0; k <= sumW; k++ {
           dp[j][k] = -INF
       }
   }
   dp[0][0] = 0
   var sumV float64
   for i := 0; i < n; i++ {
       sumV += v[i]
       // update dp backwards
       for j := i + 1; j >= 1; j-- {
           for k := sumW; k >= w[i]; k-- {
               prev := dp[j-1][k-w[i]]
               if prev > -INF {
                   val := prev + v[i]
                   if val > dp[j][k] {
                       dp[j][k] = val
                   }
               }
           }
       }
   }
   // compute answers for each j
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for j := 1; j <= n; j++ {
       ans := 0.0
       for k := 0; k <= sumW; k++ {
           if dp[j][k] < -INF/2 {
               continue
           }
           water := (dp[j][k] + sumV) / 2.0
           cap := float64(k)
           ans = max(ans, math.Min(cap, water))
       }
       fmt.Fprintf(out, "%.9f ", ans)
   }
   fmt.Fprintln(out)
}
