package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   sb := []byte(s)
   // lcp for substring comparisons
   lcp := make([][]int, n+2)
   for i := range lcp {
       lcp[i] = make([]int, n+2)
   }
   for i := n - 1; i >= 0; i-- {
       for j := n - 1; j >= 0; j-- {
           if sb[i] == sb[j] {
               lcp[i][j] = lcp[i+1][j+1] + 1
           }
       }
   }
   // dp_count and seg_count
   dp := make([][]int, n+1)
   seg := make([][]int, n+1)
   inf := int(1e9)
   for i := 0; i <= n; i++ {
       dp[i] = make([]int, n+1)
       seg[i] = make([]int, n+1)
       for j := range seg[i] {
           seg[i][j] = inf
       }
   }
   dp[0][0] = 1
   seg[0][0] = 0
   // DP
   for i := 1; i <= n; i++ {
       for l := 1; l <= i; l++ {
           if sb[i-l] == '0' {
               continue
           }
           if i-l == 0 {
               dp[i][l] = 1
               seg[i][l] = 1
           } else {
               total := 0
               minSeg := inf
               // k < l always valid
               for k := 1; k < l; k++ {
                   total = (total + dp[i-l][k]) % MOD
                   if seg[i-l][k] < minSeg {
                       minSeg = seg[i-l][k]
                   }
               }
               // k == l, compare previous and current substrings
               if l <= i-l {
                   start1 := i - l - l
                   start2 := i - l
                   // compare sb[start1:start1+l] <= sb[start2:start2+l]
                   ok := false
                   common := lcp[start1][start2]
                   if common >= l || sb[start1+common] < sb[start2+common] {
                       ok = true
                   }
                   if ok {
                       total = (total + dp[i-l][l]) % MOD
                       if seg[i-l][l] < minSeg {
                           minSeg = seg[i-l][l]
                       }
                   }
               }
               dp[i][l] = total
               if minSeg < inf {
                   seg[i][l] = minSeg + 1
               }
           }
       }
   }
   // prefix big values and powers of two
   pref := make([]*big.Int, n+1)
   pow2 := make([]*big.Int, n+1)
   pref[0] = big.NewInt(0)
   pow2[0] = big.NewInt(1)
   for i := 1; i <= n; i++ {
       pow2[i] = new(big.Int).Lsh(pow2[i-1], 1)
       pref[i] = new(big.Int).Lsh(pref[i-1], 1)
       if sb[i-1] == '1' {
           pref[i].Add(pref[i], big.NewInt(1))
       }
   }
   var answerCount int
   var minMod int64 = MOD
   hasMin := false
   var minBig big.Int
   // evaluate final
   for l := 1; l <= n; l++ {
       if dp[n][l] == 0 {
           continue
       }
       answerCount = (answerCount + dp[n][l]) % MOD
       // compute value of last segment
       vBig := new(big.Int)
       tmp := new(big.Int).Mul(pref[n-l], pow2[l])
       vBig.Sub(pref[n], tmp)
       // M = segments + v
       Mbig := new(big.Int).Add(vBig, big.NewInt(int64(seg[n][l])))
       if !hasMin || Mbig.Cmp(&minBig) < 0 {
           hasMin = true
           minBig.Set(Mbig)
           vMod := new(big.Int).Mod(vBig, big.NewInt(MOD))
           sumMod := new(big.Int).Add(vMod, big.NewInt(int64(seg[n][l])))
           sumMod.Mod(sumMod, big.NewInt(MOD))
           minMod = sumMod.Int64()
       }
   }
   // output
   fmt.Println(answerCount)
   fmt.Println(minMod)
}
