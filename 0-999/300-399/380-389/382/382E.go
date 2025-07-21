package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, K int
   if _, err := fmt.Fscan(in, &n, &K); err != nil {
       return
   }
   // Precompute factorials and inverses up to n
   maxN := n
   fac := make([]int64, maxN+1)
   ifac := make([]int64, maxN+1)
   fac[0] = 1
   for i := 1; i <= maxN; i++ {
       fac[i] = fac[i-1] * int64(i) % MOD
   }
   ifac[maxN] = modInv(fac[maxN])
   for i := maxN; i > 0; i-- {
       ifac[i-1] = ifac[i] * int64(i) % MOD
   }
   // binomial C
   C := func(n, k int) int64 {
       if k < 0 || k > n {
           return 0
       }
       return fac[n] * ifac[k] % MOD * ifac[n-k] % MOD
   }
   inv2 := modInv(2)
   // dp_count[n][g0][g1]
   // max matching g0 <= n/2
   maxM := (n + 1) / 2
   dp := make([][][]int64, n+1)
   for i := range dp {
       dp[i] = make([][]int64, maxM+2)
       for j := range dp[i] {
           dp[i][j] = make([]int64, maxM+2)
       }
   }
   // Base: one node
   dp[1][0][0] = 1
   // Build DP
   for sz := 2; sz <= n; sz++ {
       // one child (d=1)
       i := sz - 1
       // choose label set of size i from sz-1 and root label among i labels
       coeff1 := C(sz-1, i) * int64(i) % MOD
       for g0 := 0; g0 <= i/2; g0++ {
           for g1 := 0; g1 <= g0; g1++ {
               cnt := dp[i][g0][g1]
               if cnt == 0 {
                   continue
               }
               parent_g1 := g0
               mg0 := g0
               if 1+g1 > mg0 {
                   mg0 = 1 + g1
               }
               dp[sz][mg0][parent_g1] = (dp[sz][mg0][parent_g1] + cnt*coeff1) % MOD
           }
       }
       // two children (d=2)
       for i = 1; i <= sz-2; i++ {
           j := sz - 1 - i
           comb := C(sz-1, i)
           // choose label sets S1,S2 and root labels; for unordered children divide by 2
           var labelCoeff int64
           if i != j {
               labelCoeff = comb * int64(i) % MOD * int64(j) % MOD
           } else {
               labelCoeff = comb * int64(i) % MOD * int64(j) % MOD * inv2 % MOD
           }
           // iterate child1 params
           for g0_1 := 0; g0_1 <= i/2; g0_1++ {
               for g1_1 := 0; g1_1 <= g0_1; g1_1++ {
                   cnt1 := dp[i][g0_1][g1_1]
                   if cnt1 == 0 {
                       continue
                   }
                   // iterate child2 params
                   for g0_2 := 0; g0_2 <= j/2; g0_2++ {
                       for g1_2 := 0; g1_2 <= g0_2; g1_2++ {
                           cnt2 := dp[j][g0_2][g1_2]
                           if cnt2 == 0 {
                               continue
                           }
                           // pair any shapes of two subtrees
                           pairCount := cnt1 * cnt2 % MOD
                           ways := pairCount * labelCoeff % MOD
                           parent_g1 := g0_1 + g0_2
                           match1a := 1 + g1_1 + g0_2
                           match1b := 1 + g0_1 + g1_2
                           mg0 := parent_g1
                           if match1a > mg0 {
                               mg0 = match1a
                           }
                           if match1b > mg0 {
                               mg0 = match1b
                           }
                           dp[sz][mg0][parent_g1] = (dp[sz][mg0][parent_g1] + ways) % MOD
                       }
                   }
               }
           }
       }
   }
   // Sum over g1 for root free state g0 = K
   var ans int64
   if K <= maxM {
       for g1 := 0; g1 <= K; g1++ {
           ans = (ans + dp[n][K][g1]) % MOD
       }
   }
   fmt.Println(ans)
}

// modInv computes modular inverse of a mod MOD
func modInv(a int64) int64 {
   return modPow(a, MOD-2)
}

// modPow computes a^e mod MOD
func modPow(a, e int64) int64 {
   res := int64(1)
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
