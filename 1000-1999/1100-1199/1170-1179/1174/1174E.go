package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   // spf sieve
   spf := make([]int, n+1)
   spf[1] = 1
   primes := make([]int, 0, n/10)
   for i := 2; i <= n; i++ {
       if spf[i] == 0 {
           spf[i] = i
           primes = append(primes, i)
       }
       for _, p := range primes {
           if p > spf[i] || i*p > n {
               break
           }
           spf[i*p] = p
       }
   }
   // omega and find max
   omega := make([]int, n+1)
   omega[1] = 0
   M := 0
   cands := make([]int, 0)
   for i := 2; i <= n; i++ {
       omega[i] = omega[i/spf[i]] + 1
       if omega[i] > M {
           M = omega[i]
           cands = cands[:0]
           cands = append(cands, i)
       } else if omega[i] == M {
           cands = append(cands, i)
       }
   }
   // factorials and inverses
   fact := make([]int, n+1)
  fact[0] = 1
  for i := 1; i <= n; i++ {
      fact[i] = int(int64(fact[i-1]) * int64(i) % mod)
  }
   // inverses up to n
   inv := make([]int, n+1)
   inv[1] = 1
   for i := 2; i <= n; i++ {
       inv[i] = mod - int(int64(mod/i)*int64(inv[mod%i])%mod)
   }
   fullFact := fact[n-1]
   res := 0
   // process each candidate a0
   for _, a0 := range cands {
       // factor a0
       tmp := a0
       // map primes
       pm := make([]int, 0)
       ex := make([]int, 0)
       for tmp > 1 {
           p := spf[tmp]
           cnt := 0
           for tmp%p == 0 {
               tmp /= p
               cnt++
           }
           pm = append(pm, p)
           ex = append(ex, cnt)
       }
       k := len(pm)
       // dims and strides
       dims := make([]int, k)
       strides := make([]int, k)
       total := 1
       for i := 0; i < k; i++ {
           dims[i] = ex[i] + 1
       }
       for i := 0; i < k; i++ {
           strides[i] = total
           total *= dims[i]
       }
       // DP array stores sum of products of b_j/S_j
       dp := make([]int, total)
       dp[0] = 1
       // precompute base count of multiples of a0
       baseCnt := n / a0
       // iterate states
       for idx := 0; idx < total; idx++ {
           curVal := dp[idx]
           if curVal == 0 {
               continue
           }
           // decode r vector and compute sumRemoved and gPrev
           rem := idx
           sumRemoved := 0
           r := make([]int, k)
           gPrev := a0
           for i := 0; i < k; i++ {
               ri := rem % dims[i]
               rem /= dims[i]
               r[i] = ri
               sumRemoved += ri
               for t := 0; t < ri; t++ {
                   gPrev /= pm[i]
               }
           }
           if sumRemoved == M {
               continue
           }
           // compute S_j = (n-1) - (total used so far) = (n-1) - (floor(n/gPrev) - baseCnt)
           used := n/gPrev - baseCnt
           S := (n - 1) - used
           invS := inv[S]
           // transitions: remove one prime pm[i]
           for i := 0; i < k; i++ {
               if r[i] < ex[i] {
                   gNew := gPrev / pm[i]
                   b := n/gNew - n/gPrev
                   mult := int(int64(b) * int64(invS) % mod)
                   nextIdx := idx + strides[i]
                   dp[nextIdx] = (dp[nextIdx] + int(int64(curVal)*int64(mult)%mod)) % mod
               }
           }
       }
       // full removal index
       sumFrac := dp[total-1]
       // add count = sumFrac * (n-1)! 
       res = (res + int(int64(sumFrac)*int64(fullFact)%mod)) % mod
   }
   fmt.Println(res)
}
