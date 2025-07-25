package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k, m int
   if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
       return
   }
   // dp[i][M]: number of ways after i visits, current max = M
   // We only keep current and next layers
   dp := make([]int64, n+2)
   dp2 := make([]int64, n+2)
   // i = 1: starting at M = x1
   for M := 1; M <= n; M++ {
       dp[M] = 1
   }
   // steps 2..k
   for i := 1; i < k; i++ {
       // reset dp2
       for idx := 1; idx <= n; idx++ {
           dp2[idx] = 0
       }
       for M := i; M <= n; M++ {
           v := dp[M]
           if v == 0 {
               continue
           }
           // unvisited in [1..M]: M - i
           stay := int64(M - i)
           if stay > 0 {
               dp2[M] = (dp2[M] + v*stay) % MOD
           }
           // new positions y in M+1 .. min(n, M+m)
           end := M + m
           if end > n {
               end = n
           }
           for y := M + 1; y <= end; y++ {
               dp2[y] = (dp2[y] + v) % MOD
           }
       }
       // swap dp and dp2
       dp, dp2 = dp2, dp
   }
   // sum dp[M] for M>=k to n
   var ans int64
   for M := k; M <= n; M++ {
       ans = (ans + dp[M]) % MOD
   }
   fmt.Fprintln(out, ans)
}
