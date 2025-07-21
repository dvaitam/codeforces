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
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // w[i][k] weight capacity of scale at row i, position k
   w := make([][]int64, n+1)
   for i := 1; i <= n; i++ {
       cols := n - i + 1
       w[i] = make([]int64, cols+2)
       for k := 1; k <= cols; k++ {
           fmt.Fscan(reader, &w[i][k])
       }
   }
   // dp[i][k] = max total oats that can reach scale i,k and break it
   dp := make([][]int64, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = make([]int64, n+2)
   }
   // first row
   for k := 1; k <= n; k++ {
       if a[k] >= w[1][k] {
           dp[1][k] = a[k]
       }
   }
   // DP downwards
   for i := 2; i <= n; i++ {
       for k := 1; k <= n-i+1; k++ {
           var sum int64
           // from parent above left (i-1, k)
           if dp[i-1][k] > 0 {
               sum += dp[i-1][k]
           }
           // from parent above right (i-1, k+1)
           if dp[i-1][k+1] > 0 {
               sum += dp[i-1][k+1]
           }
           if sum >= w[i][k] {
               dp[i][k] = sum
           }
       }
   }
   // check bottom
   if dp[n][1] > 0 {
       fmt.Println("Cerealguy")
   } else {
       fmt.Println("Fat Rat")
   }
}
