package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 998244353

func add(a, b int64) int64 {
   c := a + b
   if c >= MOD {
       c -= MOD
   }
   return c
}
func sub(a, b int64) int64 {
   c := a - b
   if c < 0 {
       c += MOD
   }
   return c
}
func mul(a, b int64) int64 {
   return (a * b) % MOD
}
func powmod(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 != 0 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   fmt.Fscan(in, &m)
   var maxA int = 0
   // read input
   freqs := make([]int64, 100001)
   arrA := make([]int, m)
   arrC := make([]int64, m)
   for i := 0; i < m; i++ {
       var a int
       var c int64
       fmt.Fscan(in, &a, &c)
       arrA[i] = a
       arrC[i] = c
       freqs[a] += c
       if a > maxA {
           maxA = a
       }
   }
   // mobius
   mu := make([]int64, maxA+1)
   primes := make([]int, 0, maxA/10)
   isComp := make([]bool, maxA+1)
   mu[1] = 1
   for i := 2; i <= maxA; i++ {
       if !isComp[i] {
           primes = append(primes, i)
           mu[i] = -1
       }
       for _, p := range primes {
           if i*p > maxA {
               break
           }
           isComp[i*p] = true
           if i%p == 0 {
               mu[i*p] = 0
               break
           } else {
               mu[i*p] = -mu[i]
           }
       }
   }
   inv2 := mul((MOD+1)/2, 1)   // 2^{-1}
   inv4 := mul(inv2, inv2)      // 4^{-1}
   inv8 := mul(inv4, inv2)      // 8^{-1}
   var ans int64 = 0
   // precompute freq map entries
   // loop d
   for d := 1; d <= maxA; d++ {
       var sum_c1 int64
       var sum_c2 int64
       var sum_c3 int64
       var sum_bc1 int64
       var sum_bc2 int64
       var sum_bc3 int64
       var sum_b2c1 int64
       var sum_b2c2 int64
       var sum_b2c3 int64
       var expo int64
       for k := d; k <= maxA; k += d {
           c := freqs[k]
           if c == 0 {
               continue
           }
           // exponent sum
           expo = (expo + (c % (MOD-1))) % (MOD - 1)
           // c1, c2, c3
           c1 := c % MOD
           c2 := c1 * ((c1 - 1 + MOD) % MOD) % MOD
           c3 := c2 * ((c1 - 2 + MOD) % MOD) % MOD
           sum_c1 = add(sum_c1, c1)
           sum_c2 = add(sum_c2, c2)
           sum_c3 = add(sum_c3, c3)
           b := int64(k / d)
           bm := b % MOD
           // bc
           sum_bc1 = add(sum_bc1, mul(bm, c1))
           sum_bc2 = add(sum_bc2, mul(bm, c2))
           sum_bc3 = add(sum_bc3, mul(bm, c3))
           // b2c
           b2m := bm * bm % MOD
           sum_b2c1 = add(sum_b2c1, mul(b2m, c1))
           sum_b2c2 = add(sum_b2c2, mul(b2m, c2))
           sum_b2c3 = add(sum_b2c3, mul(b2m, c3))
       }
       if sum_c1 == 0 {
           continue
       }
       // base P
       P := powmod(2, expo)
       // precompute sums
       S_p1 := mul(sum_c1, inv2)
       // S_b2_alpha
       S_b2_alpha := add(add(mul(sum_b2c1, inv2), mul(mul(sum_b2c2, 3), inv4)), mul(sum_b2c3, inv8))
       // S_b2_beta
       S_b2_beta := add(mul(sum_b2c1, inv2), mul(sum_b2c2, inv4))
       // S_b2_beta_e1
       S_b2_beta_e1 := add(mul(sum_b2c1, inv4), mul(sum_b2c2, add(inv4, inv8)))
       // S_beta_b
       S_beta_b := add(mul(sum_bc1, inv2), mul(sum_bc2, inv4))
       // S_be1
       S_be1 := mul(sum_bc1, inv2)
       // S_b2_e12
       S_b2_e12 := mul(add(sum_b2c1, sum_b2c2), inv4)
       // S_be12
       S_be12 := mul(add(sum_bc1, sum_bc2), inv4)
       // S_b2_e13
       S_b2_e13 := mul(add(add(sum_b2c1, mul(sum_b2c2, 3)), sum_b2c3), inv8)
       // compute A, B, C, D
       A := sub(add(S_b2_alpha, mul(S_p1, S_b2_beta)), S_b2_beta_e1)
       B := mul(2, sub(mul(S_be1, S_beta_b), S_b2_beta_e1))
       tmp1 := sub(mul(S_be1, S_be1), S_b2_e12)
       tmp2 := sub(mul(S_be12, S_be1), S_b2_e13)
       C := sub(mul(S_p1, tmp1), mul(2, tmp2))
       D := sub(add(add(A, B), C), S_b2_beta)
       term1 := mul(P, D)
       // term2 = P^2 * inv4 * (sum_bc1^2 - (sum_b2c1+sum_b2c2))
       bracket := sub(mul(sum_bc1, sum_bc1), add(sum_b2c1, sum_b2c2))
       term2 := mul(mul(mul(P, P), inv4), bracket)
       h := sub(term1, term2)
       f := mul(mul(int64(d*d%MOD), h), mu[d])
       ans = add(ans, f)
   }
   fmt.Println((ans%MOD+MOD)%MOD)
}
