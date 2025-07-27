package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
      return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
      return
   }
   const INF = 1000000000
   for ; t > 0; t-- {
      var n int
      fmt.Fscan(in, &n)
      a := make([]int, n)
      for i := 0; i < n; i++ {
         fmt.Fscan(in, &a[i])
      }
      // dp0[i]: min skip points to clear first i bosses, next is friend's turn
      // dp1[i]: ... next is user's turn
      dp0 := make([]int, n+3)
      dp1 := make([]int, n+3)
      for i := 0; i <= n+2; i++ {
         dp0[i] = INF
         dp1[i] = INF
      }
      dp0[0] = 0
      for i := 0; i < n; i++ {
         // friend's turn at i
         if dp0[i] < INF {
            // kill one boss
            cost := dp0[i] + a[i]
            dp1[i+1] = min(dp1[i+1], cost)
            // kill two bosses
            if i+1 < n {
               cost2 := dp0[i] + a[i] + a[i+1]
               dp1[i+2] = min(dp1[i+2], cost2)
            }
         }
         // user's turn at i
         if dp1[i] < INF {
            // kill one boss
            dp0[i+1] = min(dp0[i+1], dp1[i])
            // kill two bosses
            if i+1 < n {
               dp0[i+2] = min(dp0[i+2], dp1[i])
            }
         }
      }
      res := min(dp0[n], dp1[n])
      fmt.Fprintln(out, res)
   }
}
