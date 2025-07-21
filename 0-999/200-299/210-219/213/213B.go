package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func modPow(a, e int64) int64 {
   res := int64(1)
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
   fmt.Fscan(in, &n)
   a := make([]int, 10)
   for i := 0; i < 10; i++ {
       fmt.Fscan(in, &a[i])
   }
   // sum of requirements
   sumA := 0
   for _, v := range a {
       sumA += v
   }
   // precompute factorials and invfacts up to n
   maxN := n
   // also need up to max of a[i]
   for _, v := range a {
       if v > maxN {
           maxN = v
       }
   }
   // factorials
   fact := make([]int64, maxN+1)
   invfact := make([]int64, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invfact[maxN] = modPow(fact[maxN], mod-2)
   for i := maxN; i > 0; i-- {
       invfact[i-1] = invfact[i] * int64(i) % mod
   }
   // powers of 10
   pow10 := make([]int64, n+1)
   pow10[0] = 1
   for i := 1; i <= n; i++ {
       pow10[i] = pow10[i-1] * 10 % mod
   }
   // inv_a_fact = product invfact[a[i]]
   invAFact := int64(1)
   for i := 0; i < 10; i++ {
       if a[i] > maxN {
           // impossible any sequence
           fmt.Println(0)
           return
       }
       invAFact = invAFact * invfact[a[i]] % mod
   }
   // compute a0' and inv_a_prime_fact and sumA'
   a0p := a[0]
   if a0p > 0 {
       a0p--
   }
   sumAp := sumA
   if a[0] > 0 {
       sumAp--
   }
   invAPrime := invfact[a0p]
   for i := 1; i < 10; i++ {
       invAPrime = invAPrime * invfact[a[i]] % mod
   }
   var ans int64
   // lengths from 1 to n
   for L := 1; L <= n; L++ {
       var tot int64
       if L >= sumA {
           // tot = fact[L] * 10^(L-sumA) / (L-sumA)! * invAFact
           tot = fact[L] * pow10[L-sumA] % mod * invfact[L-sumA] % mod * invAFact % mod
       }
       // invalid starting with zero
       var inv int64
       if L-1 >= sumAp {
           inv = fact[L-1] * pow10[L-1-sumAp] % mod * invfact[L-1-sumAp] % mod * invAPrime % mod
       }
       ans = (ans + tot - inv + mod) % mod
   }
   fmt.Println(ans)
}
