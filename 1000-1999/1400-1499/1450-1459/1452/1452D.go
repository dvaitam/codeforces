package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

// fast exponentiation: a^e mod mod
func modPow(a, e int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   dp := make([]int64, n+1)
   // prefix sums by parity: sum0[i] = sum dp[j] for j<=i, j%2==0
   sum0 := make([]int64, n+1)
   sum1 := make([]int64, n+1)
   dp[0] = 1
   sum0[0] = 1
   sum1[0] = 0
   for i := 1; i <= n; i++ {
       // dp[i] = sum of dp[j] for j<=i-1, (i-j)%2==1 => j%2 == (i-1)%2
       if (i-1)&1 == 0 {
           dp[i] = sum0[i-1]
       } else {
           dp[i] = sum1[i-1]
       }
       if dp[i] >= mod {
           dp[i] %= mod
       }
       // update prefix sums
       if i&1 == 0 {
           sum0[i] = (sum0[i-1] + dp[i]) % mod
           sum1[i] = sum1[i-1]
       } else {
           sum1[i] = (sum1[i-1] + dp[i]) % mod
           sum0[i] = sum0[i-1]
       }
   }
   // total subsets is 2^n, inverse
   total := modPow(2, int64(n))
   invTotal := modPow(total, mod-2)
   ans := dp[n] * invTotal % mod
   fmt.Fprintln(writer, ans)
}
