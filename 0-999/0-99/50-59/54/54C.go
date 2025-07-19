package main

import "fmt"

// getn returns the count of numbers x in [0, a] whose most significant digit is 1.
func getn(a int64) int64 {
   if a <= 0 {
       return 0
   }
   var c int64 = 1
   // find highest power of 10 <= a
   for i := 0; i < 20; i++ {
       if a/c == 0 {
           c /= 10
           break
       }
       c *= 10
   }
   var ans int64
   // count numbers with msd=1 at highest digit
   if a/c != 1 {
       ans += c
   } else {
       ans += a%c + 1
   }
   // count full blocks for lower digit positions
   c /= 10
   for c > 0 {
       ans += c
       c /= 10
   }
   return ans
}

func main() {
   var n int
   fmt.Scan(&n)
   pro := make([]float64, n)
   for i := 0; i < n; i++ {
       var l, r int64
       fmt.Scan(&l, &r)
       total := r - l + 1
       cnt := getn(r) - getn(l-1)
       pro[i] = float64(cnt) / float64(total)
   }
   var kPercent int
   fmt.Scan(&kPercent)
   // minimal number of variables with first digit 1
   k := (n*kPercent + 99) / 100
   // dp[i][j]: probability that among first i variables, exactly j have msd=1
   dp := make([][]float64, n+1)
   for i := range dp {
       dp[i] = make([]float64, k+1)
   }
   dp[0][0] = 1.0
   for i := 1; i <= n; i++ {
       for j := 0; j <= i && j <= k; j++ {
           if j == 0 {
               dp[i][j] = dp[i-1][j] * (1 - pro[i-1])
           } else {
               dp[i][j] = dp[i-1][j-1]*pro[i-1] + dp[i-1][j]*(1-pro[i-1])
           }
       }
   }
   var sum float64
   for j := 0; j < k; j++ {
       sum += dp[n][j]
   }
   ans := 1.0 - sum
   if ans < 0 {
       ans = 0
   } else if ans > 1 {
       ans = 1
   }
   fmt.Printf("%.12f\n", ans)
}
