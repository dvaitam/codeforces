package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([][]int64, n+1)
   for i := 1; i <= n; i++ {
       a[i] = make([]int64, n+1)
       for j := 1; j <= n; j++ {
           fmt.Fscan(in, &a[i][j])
       }
   }
   const INF = int64(4e18)
   // best 2-cycle: 1 -> i -> 1
   best2 := INF
   for i := 2; i <= n; i++ {
       cost := a[1][i] + a[i][1]
       if cost < best2 {
           best2 = cost
       }
   }
   // best 4-cycle: 1 -> i -> j -> i -> 1 (i,j != 1, i != j)
   best4 := INF
   for i := 2; i <= n; i++ {
       for j := 2; j <= n; j++ {
           if i == j {
               continue
           }
           cost := a[1][i] + a[i][j] + a[j][i] + a[i][1]
           if cost < best4 {
               best4 = cost
           }
       }
   }
   // DP on even lengths up to k
   dp := make([]int64, k+1)
   for i := 1; i <= k; i++ {
       dp[i] = INF
   }
   dp[0] = 0
   for length := 2; length <= k; length += 2 {
       // use one 2-cycle
       if dp[length-2] != INF {
           v := dp[length-2] + best2
           if v < dp[length] {
               dp[length] = v
           }
       }
       // use one 4-cycle
       if length >= 4 && dp[length-4] != INF {
           v := dp[length-4] + best4
           if v < dp[length] {
               dp[length] = v
           }
       }
   }
   // output answer
   fmt.Println(dp[k])
}
