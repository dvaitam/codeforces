package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, d int
   if _, err := fmt.Fscan(reader, &n, &d); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   sort.Ints(c)
   // prefix sums
   prefix := make([]int, n+1)
   for i := 1; i <= n; i++ {
       prefix[i] = prefix[i-1] + c[i-1]
   }
   const INF = 1e9
   dp := make([]int, n+1)
   for i := range dp {
       dp[i] = INF
   }
   dp[0] = 0
   // dp[i]: min days to obtain prefix items sum = prefix[i]
   for i := 1; i <= n; i++ {
       for j := 0; j < i; j++ {
           if dp[j] < INF && 2*prefix[j] + d >= prefix[i] {
               if dp[j]+1 < dp[i] {
                   dp[i] = dp[j] + 1
               }
           }
       }
   }
   // find best prefix
   bestSum, bestDays := 0, 0
   for i := 0; i <= n; i++ {
       if dp[i] < INF {
           if prefix[i] > bestSum {
               bestSum = prefix[i]
               bestDays = dp[i]
           } else if prefix[i] == bestSum && dp[i] < bestDays {
               bestDays = dp[i]
           }
       }
   }
   fmt.Fprintf(writer, "%d %d", bestSum, bestDays)
}
