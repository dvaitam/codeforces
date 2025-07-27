package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// modpow computes a^e mod MOD
func modpow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

// modinv computes modular inverse of a under MOD
func modinv(a int64) int64 {
   // MOD is prime
   return modpow((a%MOD+MOD)%MOD, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int64
   var k, m int
   if _, err := fmt.Fscan(in, &n, &k, &m); err != nil {
       return
   }
   // Precompute answers S[N] for N = 0..k by direct DP
   K := k
   S := make([]int64, K+1)
   for N := 0; N <= K; N++ {
       if N < k {
           S[N] = 0
           continue
       }
       // dp arrays of size N+2
       dp := make([]int64, N+2)
       dp2 := make([]int64, N+2)
       // initial: dp[1][M] = 1
       for M := 1; M <= N; M++ {
           dp[M] = 1
       }
       // steps 2..k
       for step := 1; step < k; step++ {
           for i := 1; i <= N; i++ {
               dp2[i] = 0
           }
           for M := step; M <= N; M++ {
               v := dp[M]
               if v == 0 {
                   continue
               }
               // stay within existing max
               stay := int64(M - step)
               if stay > 0 {
                   dp2[M] = (dp2[M] + v*stay) % MOD
               }
               // move to new max in (M+1 .. M+m)
               end := M + m
               if end > N {
                   end = N
               }
               for y := M + 1; y <= end; y++ {
                   dp2[y] = (dp2[y] + v) % MOD
               }
           }
           dp, dp2 = dp2, dp
       }
       // sum dp[k][M] for M=k..N
       var sum int64
       for M := k; M <= N; M++ {
           sum = (sum + dp[M]) % MOD
       }
       S[N] = sum
   }
   // if n small, output directly
   if n <= int64(K) {
       fmt.Fprintln(out, S[n])
       return
   }
   // Lagrange interpolation on points (x=i, y=S[i]) for i=0..K
   // Precompute factorials and invfactorials
   fact := make([]int64, K+1)
   invfact := make([]int64, K+1)
   fact[0] = 1
   for i := 1; i <= K; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   invfact[K] = modinv(fact[K])
   for i := K; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % MOD
   }
   // prefix and suffix products of (n - j)
   pref := make([]int64, K+1)
   suf := make([]int64, K+1)
   pref[0] = 1
   for i := 1; i <= K; i++ {
       pref[i] = pref[i-1] * ((n - int64(i-1) + MOD) % MOD) % MOD
   }
   suf[K] = 1
   for i := K - 1; i >= 0; i-- {
       suf[i] = suf[i+1] * ((n - int64(i+1) + MOD) % MOD) % MOD
   }
   // compute interpolation
   var ans int64
   for i := 0; i <= K; i++ {
       // numerator = ∏_{j≠i}(n - j)
       num := pref[i] * suf[i] % MOD
       // denom = ∏_{j≠i}(i - j) = (-1)^{K-i} * i! * (K-i)!
       sign := int64(1)
       if (K-i)&1 == 1 {
           sign = MOD - 1
       }
       invDen := invfact[i] * invfact[K-i] % MOD * sign % MOD
       ans = (ans + S[i] * num % MOD * invDen) % MOD
   }
   fmt.Fprintln(out, ans)
}
