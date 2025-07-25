package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func add(a, b int) int {
   a += b
   if a >= mod {
       a -= mod
   }
   return a
}
func sub(a, b int) int {
   a -= b
   if a < 0 {
       a += mod
   }
   return a
}
func mul(a, b int) int {
   return int((int64(a) * int64(b)) % mod)
}
func powmod(a, e int) int {
   res := 1
   x := a
   for e > 0 {
       if e&1 != 0 {
           res = mul(res, x)
       }
       x = mul(x, x)
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   var l int64
   if _, err := fmt.Fscan(in, &n, &k, &l); err != nil {
       return
   }
   // precompute factorials
   maxN := 2*n + 1
   fact := make([]int, maxN+1)
   invf := make([]int, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = mul(fact[i-1], i)
   }
   invf[maxN] = powmod(fact[maxN], mod-2)
   for i := maxN; i > 0; i-- {
       invf[i-1] = mul(invf[i], i)
   }
   // pow2
   pow2 := make([]int, n+1)
   pow2[0] = 1
   for i := 1; i <= n; i++ {
       pow2[i] = add(pow2[i-1], pow2[i-1])
   }
   // C(n,r)
   cnr := make([]int, n+1)
   for r := 0; r <= n; r++ {
       cnr[r] = mul(fact[n], mul(invf[r], invf[n-r]))
   }
   // F[r] = ∫0^1 (2t(1-t))^r dt = 2^r * r! * r! / (2r+1)!
   Fr := make([]int, n+1)
   for r := 0; r <= n; r++ {
       // 2^r * fact[r] * fact[r] * invf[2r+1]
       Fr[r] = mul(pow2[r], mul(fact[r], mul(fact[r], invf[2*r+1])))
   }
   // compute sumS = sum_{i=0..k-1} S_i
   sumS := 0
   for i := 0; i < k; i++ {
       // S_i = sum_{r=i..n} C(n,r)*C(r,i)*(-1)^{r-i} * Fr[r]
       // C(r,i) = r!/(i!(r-i)!) = mul(fact[r], mul(invf[i], invf[r-i]))
       for r := i; r <= n; r++ {
           coef := mul(cnr[r], mul(fact[r], mul(invf[i], invf[r-i])))
           term := mul(coef, Fr[r])
           if (r-i)&1 == 1 {
               term = mod - term
           }
           sumS = add(sumS, term)
       }
   }
   // result = l * (1 - sumS)
   res := int((l % mod) * int64(sub(1, sumS)) % mod)
   if res < 0 {
       res += mod
   }
   fmt.Println(res)
}
