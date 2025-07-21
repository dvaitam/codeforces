package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, h int
   if _, err := fmt.Fscan(reader, &n, &h); err != nil {
       return
   }
   // dp[i][j]: number of BSTs with i nodes and height <= j
   dp := make([][]uint64, n+1)
   for i := range dp {
       dp[i] = make([]uint64, n+1)
   }
   // empty tree has height 0
   for j := 0; j <= n; j++ {
       dp[0][j] = 1
   }
   // build DP
   for j := 1; j <= n; j++ {
       for i := 1; i <= n; i++ {
           var cnt uint64
           for left := 0; left < i; left++ {
               right := i - 1 - left
               cnt += dp[left][j-1] * dp[right][j-1]
           }
           dp[i][j] = cnt
       }
   }
   // total number of BSTs of size n is dp[n][n]
   total := dp[n][n]
   // number with height < h is dp[n][h-1]
   less := uint64(0)
   if h > 0 {
       less = dp[n][h-1]
   }
   // answer: number with height >= h
   ans := total - less
   fmt.Println(ans)
}
