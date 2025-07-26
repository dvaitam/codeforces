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
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // can[i][j] = value if segment [i..j] can shrink to one element, else 0
   can := make([][]int, n)
   dp := make([][]int, n)
   for i := 0; i < n; i++ {
       can[i] = make([]int, n)
       dp[i] = make([]int, n)
       can[i][i] = a[i]
       dp[i][i] = 1
   }
   const INF = 1e9
   // DP over lengths
   for length := 2; length <= n; length++ {
       for l := 0; l+length-1 < n; l++ {
           r := l + length - 1
           dp[l][r] = INF
           // compute can[l][r]
           for k := l; k < r; k++ {
               v1 := can[l][k]
               if v1 == 0 {
                   continue
               }
               v2 := can[k+1][r]
               if v2 == v1 {
                   can[l][r] = v1 + 1
                   break
               }
           }
           // if fully shrinkable
           if can[l][r] != 0 {
               dp[l][r] = 1
               continue
           }
           // else split
           for k := l; k < r; k++ {
               dp[l][r] = min(dp[l][r], dp[l][k] + dp[k+1][r])
           }
       }
   }
   fmt.Fprintln(writer, dp[0][n-1])
}
