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
   // dp table of size (n+1)x(n+1)
   dp := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int, n+1)
   }
   // initialize first row and first column to 1
   for i := 1; i <= n; i++ {
       dp[i][1] = 1
       dp[1][i] = 1
   }
   // fill dp by sum of top and left
   for i := 2; i <= n; i++ {
       for j := 2; j <= n; j++ {
           dp[i][j] = dp[i-1][j] + dp[i][j-1]
       }
   }
   // the maximum value is at dp[n][n]
   fmt.Fprintln(writer, dp[n][n])
}
