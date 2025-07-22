package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   const mod = 1000000007
   dp := make([]int64, n+1)
   // dp[i] = moves to reach room i+1 from room 1
   dp[0] = 0
   for i := 1; i <= n; i++ {
       // dp[i] = 2 + 2*dp[i-1] - dp[p[i]-1]
       val := 2 + 2*dp[i-1] - dp[p[i]-1]
       // mod adjust
       val %= mod
       if val < 0 {
           val += mod
       }
       dp[i] = val
   }
   fmt.Println(dp[n])
}
