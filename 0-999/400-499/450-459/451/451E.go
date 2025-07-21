package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const mod = 1000000007

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
   in := bufio.NewReader(os.Stdin)
   var n int
   var s int64
   if _, err := fmt.Fscan(in, &n, &s); err != nil {
       return
   }
   fplus := make([]int64, n)
   for i := 0; i < n; i++ {
       var fi int64
       fmt.Fscan(in, &fi)
       fplus[i] = fi + 1
   }
   // Precompute factorials and inverse factorials up to n
   maxR := n
   fact := make([]int64, maxR+1)
   invFact := make([]int64, maxR+1)
   fact[0] = 1
   for i := 1; i <= maxR; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invFact[maxR] = modPow(fact[maxR], mod-2)
   for i := maxR; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % mod
   }
   // Prepare subset sums and popcounts
   N := 1 << n
   sumF := make([]int64, N)
   pc := make([]int, N)
   for mask := 1; mask < N; mask++ {
       lb := mask & -mask
       i := bits.TrailingZeros(uint(lb))
       prev := mask ^ lb
       sumF[mask] = sumF[prev] + fplus[i]
       pc[mask] = pc[prev] + 1
   }
   // inclusion-exclusion
   var ans int64
   r := n - 1
   for mask := 0; mask < N; mask++ {
       k := pc[mask]
       total := sumF[mask]
       a := s - total + int64(n-1)
       if a < int64(r) {
           continue
       }
       // compute C(a, r)
       var num int64 = 1
       for i := 0; i < r; i++ {
           t := (a - int64(i)) % mod
           if t < 0 {
               t += mod
           }
           num = num * t % mod
       }
       comb := num * invFact[r] % mod
       if k&1 == 1 {
           ans = (ans - comb + mod) % mod
       } else {
           ans = (ans + comb) % mod
       }
   }
   // if n == 0? but n>=1 by constraints
   fmt.Println(ans)
}
