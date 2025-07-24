package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix sums
   pref := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       pref[i] = pref[i-1] + a[i]
   }
   // dp[j]: optimal score difference for state at prefix ending j
   dp := make([]int64, n+2)
   // mx[j]: max over k>=j of (pref[k] - dp[k])
   mx := make([]int64, n+3)
   const inf = int64(1) << 60
   // initialize
   dp[n] = 0
   mx[n+1] = -inf
   // compute mx[n]
   mx[n] = pref[n] - dp[n]
   // fill dp and mx backwards
   for j := n - 1; j >= 1; j-- {
       // best move: choose k > j to maximize pref[k] - dp[k]
       dp[j] = mx[j+1]
       // update mx for j
       v := pref[j] - dp[j]
       // mx[j] = max(v, mx[j+1])
       if v > mx[j+1] {
           mx[j] = v
       } else {
           mx[j] = mx[j+1]
       }
   }
   // answer is dp[1]
   fmt.Fprintln(writer, dp[1])
}
