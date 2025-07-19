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

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n+3)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   k := (n + 1) / 2
   const INF = 1000000000000000000
   // dp[2][k+2][3]
   dp := make([][][]int, 2)
   for t := 0; t < 2; t++ {
       dp[t] = make([][]int, k+2)
       for j := 0; j < k+2; j++ {
           dp[t][j] = make([]int, 3)
           for s := 0; s < 3; s++ {
               dp[t][j][s] = INF
           }
       }
   }
   t0, t1 := 0, 1
   dp[t1][0][0] = 0
   dp[t1][1][2] = 0
   for i := 2; i <= n; i++ {
       t0 ^= 1
       t1 ^= 1
       // reset new layer
       for j := 0; j <= k; j++ {
           dp[t1][j][0] = INF
           dp[t1][j][1] = INF
           dp[t1][j][2] = INF
       }
       // Calculate adjustments
       A := max(0, a[i]-a[i-1]+1)
       B := max(0, a[i-1]-a[i]+1)
       C := min(a[i-2]-1, a[i-1]) - a[i] + 1
       if C < 0 {
           C = 0
       }
       for j := 0; j <= k; j++ {
           v0 := dp[t0][j][0]
           v1 := dp[t0][j][1]
           // state 0: normal
           if v0 < v1 {
               dp[t1][j][0] = v0
           } else {
               dp[t1][j][0] = v1
           }
           // state 1: previous was peak
           dp[t1][j][1] = dp[t0][j][2] + A
           // state 2: current is peak
           if j+1 <= k {
               x := v0 + B
               y := v1 + C
               if x < y {
                   dp[t1][j+1][2] = x
               } else {
                   dp[t1][j+1][2] = y
               }
           }
       }
   }
   // output results
   for i := 1; i <= k; i++ {
       res := dp[t1][i][0]
       if dp[t1][i][1] < res {
           res = dp[t1][i][1]
       }
       if dp[t1][i][2] < res {
           res = dp[t1][i][2]
       }
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, res)
   }
}
