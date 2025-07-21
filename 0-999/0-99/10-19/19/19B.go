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
   // dp[j]: min cost to achieve total weight j (or j>=n)
   const INF64 = 1 << 60
   dp := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = INF64
   }
   // Process items
   for i := 0; i < n; i++ {
       var ti, ci int
       fmt.Fscan(reader, &ti, &ci)
       w := ti + 1
       cost := int64(ci)
       // 0-1 knapsack backward
       for j := n; j >= 0; j-- {
           if dp[j] == INF64 {
               continue
           }
           nj := j + w
           if nj > n {
               nj = n
           }
           if dp[j]+cost < dp[nj] {
               dp[nj] = dp[j] + cost
           }
       }
   }
   // Answer is dp[n]
   fmt.Println(dp[n])
