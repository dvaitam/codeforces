package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func modPow(a, e int) int {
   res := 1
   base := a % MOD
   for e > 0 {
       if e&1 == 1 {
           res = int((int64(res) * base) % MOD)
       }
       base = int((int64(base) * base) % MOD)
       e >>= 1
   }
   return res
}

func modInv(a int) int {
   // MOD is prime
   return modPow((a%MOD+MOD)%MOD, MOD-2)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   w := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i])
   }
   // compute initial sums
   var S_plus, S_minus int
   for i := 0; i < n; i++ {
       if a[i] == 1 {
           S_plus += w[i]
       } else {
           S_minus += w[i]
       }
   }
   S := S_plus + S_minus
   // precompute inverses for possible total weights
   // S_total ranges from s_min to s_max
   sMin := S - m
   if sMin < S_plus {
       sMin = S_plus
   }
   sMax := S + m
   invSize := sMax - sMin + 1
   inv := make([]int, invSize)
   for i := 0; i < invSize; i++ {
       tot := sMin + i
       inv[i] = modInv(tot)
   }
   // dp[k][l] -> use rolling arrays dp[l]
   dp := make([]int, m+1)
   ndp := make([]int, m+1)
   dp[0] = 1
   for k := 0; k < m; k++ {
       // reset ndp
       for i := 0; i <= m; i++ {
           ndp[i] = 0
       }
       for l := 0; l <= k; l++ {
           v := dp[l]
           if v == 0 {
               continue
           }
           // total weight at this state
           S_total := S + 2*l - k
           // index for inv
           idx := S_total - sMin
           invTot := inv[idx]
           // liked transition
           pPlus := int((int64(S_plus+l) * invTot) % MOD)
           // add to l+1
           ndp[l+1] = (ndp[l+1] + int((int64(v) * pPlus) % MOD)) % MOD
           // disliked sum
           rem := S_minus - (k - l)
           if rem > 0 {
               pMinus := int((int64(rem) * invTot) % MOD)
               ndp[l] = (ndp[l] + int((int64(v) * pMinus) % MOD)) % MOD
           }
       }
       // swap dp and ndp
       dp, ndp = ndp, dp
   }
   // compute E[X]
   EX := 0
   for l := 0; l <= m; l++ {
       EX = (EX + int((int64(dp[l]) * l) % MOD)) % MOD
   }
   // compute multipliers
   invSp := modInv(S_plus)
   tPlus := int((int64(S_plus+EX) * invSp) % MOD)
   // compute m-EX
   mMinusEX := (m - EX) % MOD
   if mMinusEX < 0 {
       mMinusEX += MOD
   }
   invSm := 1
   if S_minus > 0 {
       invSm = modInv(S_minus)
   }
   tMinus := 0
   if S_minus > 0 {
       tMinus = int((int64((S_minus%MOD - mMinusEX + MOD) % MOD) * invSm) % MOD)
   }
   // output results
   for i := 0; i < n; i++ {
       var res int
       if a[i] == 1 {
           res = int((int64(w[i]) * tPlus) % MOD)
       } else {
           res = int((int64(w[i]) * tMinus) % MOD)
       }
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, res)
   }
   fmt.Fprintln(out)
}
