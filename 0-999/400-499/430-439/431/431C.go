package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k, d int
   if _, err := fmt.Fscan(in, &n, &k, &d); err != nil {
       return
   }
   // dp[s][0]: ways to get sum s without using any term >= d
   // dp[s][1]: ways to get sum s with at least one term >= d
   dp := make([][2]int, n+1)
   dp[0][0] = 1
   for s := 1; s <= n; s++ {
       var without, with int
       for i := 1; i <= k; i++ {
           if s-i < 0 {
               break
           }
           if i < d {
               // adding a small term does not change "with" status
               without = (without + dp[s-i][0]) % mod
               with = (with + dp[s-i][1]) % mod
           } else {
               // adding a large term makes it "with"
               with = (with + dp[s-i][0] + dp[s-i][1]) % mod
           }
       }
       dp[s][0] = without
       dp[s][1] = with
   }
   fmt.Println(dp[n][1])
}
