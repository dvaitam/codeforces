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

   var n, m, c0, d0 int
   fmt.Fscan(reader, &n, &m, &c0, &d0)

   dp := make([]int, n+1)
   // initialize unbounded knapsack for item with cost c0 and value d0
   for i := c0; i <= n; i++ {
       if dp[i-c0]+d0 > dp[i] {
           dp[i] = dp[i-c0] + d0
       }
   }

   // process m groups of items
   for i := 0; i < m; i++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       count := a / b
       // bounded knapsack: count items each cost c and value d
       for j := 0; j < count; j++ {
           for k := n; k >= c; k-- {
               if dp[k-c]+d > dp[k] {
                   dp[k] = dp[k-c] + d
               }
           }
       }
   }

   fmt.Fprintln(writer, dp[n])
}
