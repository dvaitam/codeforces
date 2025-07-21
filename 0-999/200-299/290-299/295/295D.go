package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modpow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 != 0 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   if n <= 0 || m < 2 {
       fmt.Fprintln(writer, 0)
       return
   }
   // Precompute f[w] = number of ways to place interval of width w: m-1-w, w=0..m-2
   maxW := m - 2
   f := make([]int, maxW+1)
   invf := make([]int64, maxW+1)
   for w := 0; w <= maxW; w++ {
       val := m - 1 - w
       f[w] = val % MOD
       invf[w] = modpow(int64(val), MOD-2)
   }
   // dp[u][w]: sum of products for non-decreasing sequences length u ending at w
   dp := make([][]int, n+1)
   dp[1] = make([]int, maxW+1)
   for w := 0; w <= maxW; w++ {
       dp[1][w] = f[w]
   }
   // Build dp for u=2..n
   for u := 2; u <= n; u++ {
       dp[u] = make([]int, maxW+1)
       sum := 0
       for w := 0; w <= maxW; w++ {
           sum = (sum + dp[u-1][w]) % MOD
           dp[u][w] = int(int64(f[w]) * int64(sum) % MOD)
       }
   }
   // Compute answer
   var ans int64
   // Temporary prefix arrays
   prefix1 := make([]int64, n+2)
   prefix2 := make([]int64, n+2)
   for w := 0; w <= maxW; w++ {
       // A[u] = dp[u][w] for u=1..n
       // Build prefix sums over u
       prefix1[0] = 0
       prefix2[0] = 0
       for u := 1; u <= n; u++ {
           av := int64(dp[u][w])
           prefix1[u] = (prefix1[u-1] + av) % MOD
           prefix2[u] = (prefix2[u-1] + av*int64(u)%MOD) % MOD
       }
       // Sum over u and v
       var sumUV int64
       for u := 1; u <= n; u++ {
           // v from 1 to n-u+1
           lim := n - u + 1
           if lim < 1 {
               break
           }
           total1 := prefix1[lim]
           total2 := prefix2[lim]
           coeff := (int64(n+2-u)*total1 - total2) % MOD
           if coeff < 0 {
               coeff += MOD
           }
           sumUV = (sumUV + int64(dp[u][w]) * coeff) % MOD
       }
       // divide by f[w] to correct double count of peak
       ans = (ans + sumUV*invf[w]) % MOD
   }
   fmt.Fprintln(writer, ans)
}
