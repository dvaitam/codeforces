package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modpow(a int64, e int64) int64 {
   var res int64 = 1
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m, k int64
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // x[1..n] periodic sequence with sum k, each 0<=x<=n
   // occurrences per position: q or q+1 times
   q := m / int64(n)
   r := int(m % int64(n))
   // precompute C(n, x)
   fac := make([]int64, n+1)
   ifac := make([]int64, n+1)
   fac[0] = 1
   for i := 1; i <= n; i++ {
       fac[i] = fac[i-1] * int64(i) % MOD
   }
   ifac[n] = modpow(fac[n], MOD-2)
   for i := n; i > 0; i-- {
       ifac[i-1] = ifac[i] * int64(i) % MOD
   }
   C := make([]int64, n+1)
   for i := 0; i <= n; i++ {
       C[i] = fac[n] * ifac[i] % MOD * ifac[n-i] % MOD
   }
   // precompute powers
   pw1 := make([]int64, n+1) // C^q
   pw2 := make([]int64, n+1) // C^(q+1)
   for x := 0; x <= n; x++ {
       pw1[x] = modpow(C[x], q)
       pw2[x] = pw1[x] * C[x] % MOD
   }
   K := int(k)
   // dp[i][s]: using first i positions sum s
   dp := make([]int64, K+1)
   dp[0] = 1
   for i := 1; i <= n; i++ {
       pw := pw1
       if i <= r {
           pw = pw2
       }
       next := make([]int64, K+1)
       // transition
       for s := 0; s <= K; s++ {
           v := dp[s]
           if v == 0 {
               continue
           }
           // choose x in 0..n, and s+x <= K
           maxx := n
           if s+maxx > K {
               maxx = K - s
           }
           for x := 0; x <= maxx; x++ {
               next[s+x] = (next[s+x] + v*pw[x]) % MOD
           }
       }
       dp = next
   }
   // result dp[k]
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprint(writer, dp[K])
}
