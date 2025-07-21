package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t [3][3]int64
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           fmt.Fscan(reader, &t[i][j])
       }
   }
   var n int
   fmt.Fscan(reader, &n)
   // dp[d][i][j]: minimal cost to move d disks from rod i to j
   var dp [41][3][3]int64
   // Base case for one disk
   for i := 0; i < 3; i++ {
       for j := 0; j < 3; j++ {
           if i == j {
               dp[1][i][j] = 0
           } else {
               k := 3 - i - j
               dp[1][i][j] = min(t[i][j], t[i][k]+t[k][j])
           }
       }
   }
   // Build up for more disks
   for d := 2; d <= n; d++ {
       for i := 0; i < 3; i++ {
           for j := 0; j < 3; j++ {
               if i == j {
                   dp[d][i][j] = 0
                   continue
               }
               k := 3 - i - j
               // classical move: via direct move of largest disk
               cost1 := dp[d-1][i][k] + t[i][j] + dp[d-1][k][j]
               // alternative: move largest via k with extra moves
               cost2 := dp[d-1][i][j] + t[i][k] + dp[d-1][j][i] + t[k][j] + dp[d-1][i][j]
               dp[d][i][j] = min(cost1, cost2)
           }
       }
   }
   // Output minimal cost to move n disks from rod 1 (0) to rod 3 (2)
   fmt.Println(dp[n][0][2])
}
