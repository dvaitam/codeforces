package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func computeS(n, m, x int, skipX bool) int64 {
   // dp[i][r]: number of ways to pick i intervals with last r = r
   dp := make([][]int64, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int64, m+1)
   }
   dp[0][0] = 1
   // temporary dp2 and prefix sum
   for l := 1; l <= m; l++ {
       // copy dp to dp2
       dp2 := make([][]int64, n+1)
       for i := 0; i <= n; i++ {
           row := make([]int64, m+1)
           copy(row, dp[i])
           dp2[i] = row
       }
       if !(skipX && l == x) {
           // we can choose interval at this l
           for i := 0; i < n; i++ {
               // prefix sums of dp[i]
               pre := make([]int64, m+1)
               pre[0] = dp[i][0]
               for r := 1; r <= m; r++ {
                   pre[r] = pre[r-1] + dp[i][r]
                   if pre[r] >= mod {
                       pre[r] -= mod
                   }
               }
               // for possible r values
               for r := l; r <= m; r++ {
                   // sum dp[i][0..r-1] = pre[r-1]
                   dp2[i+1][r] = (dp2[i+1][r] + pre[r-1]) % mod
               }
           }
       }
       dp = dp2
   }
   var total int64
   for r := 1; r <= m; r++ {
       total = (total + dp[n][r]) % mod
   }
   return total
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, x int
   if _, err := fmt.Fscan(in, &n, &m, &x); err != nil {
       return
   }
   // precompute factorial
   fact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   total := computeS(n, m, x, false)
   bad := computeS(n, m, x, true)
   ways := (total - bad + mod) % mod
   ans := ways * fact[n] % mod
   fmt.Println(ans)
}
