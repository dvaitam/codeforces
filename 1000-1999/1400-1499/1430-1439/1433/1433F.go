package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // dp2[r] = max sum so far with sum mod k == r
   const INF = 1_000_000_000
   dp2 := make([]int, k)
   for i := 1; i < k; i++ {
       dp2[i] = -INF
   }
   half := m / 2
   for i := 0; i < n; i++ {
       a := make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(in, &a[j])
       }
       // row dp: dp[c][r]
       // use 2D slice
       dp := make([][]int, half+1)
       for c := 0; c <= half; c++ {
           dp[c] = make([]int, k)
           for r := 0; r < k; r++ {
               dp[c][r] = -INF
           }
       }
       dp[0][0] = 0
       for _, v := range a {
           rem := v % k
           // update from high c to low
           for c := half - 1; c >= 0; c-- {
               for r := 0; r < k; r++ {
                   prev := dp[c][r]
                   if prev < 0 {
                       continue
                   }
                   nr := r + rem
                   if nr >= k {
                       nr -= k
                   }
                   dp[c+1][nr] = max(dp[c+1][nr], prev+v)
               }
           }
       }
       // best for this row by remainder
       rowBest := make([]int, k)
       for r := 0; r < k; r++ {
           rowBest[r] = -INF
       }
       for c := 0; c <= half; c++ {
           for r := 0; r < k; r++ {
               rowBest[r] = max(rowBest[r], dp[c][r])
           }
       }
       // combine dp2 and rowBest -> newDp
       newDp := make([]int, k)
       for i := 0; i < k; i++ {
           newDp[i] = -INF
       }
       for r1 := 0; r1 < k; r1++ {
           if dp2[r1] < 0 {
               continue
           }
           for r2 := 0; r2 < k; r2++ {
               if rowBest[r2] < 0 {
                   continue
               }
               nr := r1 + r2
               if nr >= k {
                   nr -= k
               }
               newDp[nr] = max(newDp[nr], dp2[r1]+rowBest[r2])
           }
       }
       dp2 = newDp
   }
   // answer is dp2[0]
   res := dp2[0]
   if res < 0 {
       res = 0
   }
   fmt.Fprintln(out, res)
}
