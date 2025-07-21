package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   h := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &h[i])
   }
   // prefix sums
   p := make([]int64, n+1)
   for i := 0; i < n; i++ {
       p[i+1] = p[i] + h[i]
   }
   // dp[i]: max segments for prefix up to i
   dp := make([]int, n+1)
   // bestSum[i]: minimal last segment sum for dp[i] segments
   const inf = int64(9e18)
   bestSum := make([]int64, n+1)
   dp[0] = 0
   bestSum[0] = 0
   for i := 1; i <= n; i++ {
       dp[i] = 0
       bestSum[i] = inf
       for j := i - 1; j >= 0; j-- {
           s := p[i] - p[j]
           if s < bestSum[j] {
               continue
           }
           seg := dp[j] + 1
           if seg > dp[i] || (seg == dp[i] && s < bestSum[i]) {
               dp[i] = seg
               bestSum[i] = s
           }
       }
   }
   // minimum operations = total towers - max segments
   res := n - dp[n]
   fmt.Println(res)
}
