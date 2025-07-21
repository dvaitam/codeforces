package main

import "fmt"

func max(x, y int) int {
   if x > y {
       return x
   }
   return y
}

func main() {
   var n, a, b, c int
   if _, err := fmt.Scan(&n, &a, &b, &c); err != nil {
       return
   }
   const negInf = -1000000
   dp := make([]int, n+1)
   // Initialize dp values to negative infinity except dp[0]
   for i := 1; i <= n; i++ {
       dp[i] = negInf
   }
   dp[0] = 0
   // Build up dp
   for i := 1; i <= n; i++ {
       if i >= a {
           dp[i] = max(dp[i], dp[i-a]+1)
       }
       if i >= b {
           dp[i] = max(dp[i], dp[i-b]+1)
       }
       if i >= c {
           dp[i] = max(dp[i], dp[i-c]+1)
       }
   }
   fmt.Println(dp[n])
}
