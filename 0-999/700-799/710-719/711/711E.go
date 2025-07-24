package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func modPow(a, e, mod int64) int64 {
   res := int64(1)
   a %= mod
   for e > 0 {
       if e&1 != 0 {
           res = (res * a) % mod
       }
       a = (a * a) % mod
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k uint64
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // Modulus
   const M int64 = 1000003
   phi := M - 1
   // If k > 2^n, collision guaranteed
   if n <= 63 && k > (uint64(1) << n) {
       fmt.Println(1, 1)
       return
   }
   // Compute v2(P) = n + (k-1) - popcount(k-1)
   pc := bits.OnesCount64(k - 1)
   // v2P may be large
   // Compute modular equivalents
   // n_mod and k_mod for exponent mod phi
   n_mod := int64(n % uint64(phi))
   k_mod := int64(k % uint64(phi))
   // Exponent for 2^(n*k)
   exp_nk := (n_mod * k_mod) % phi
   pow2nk := modPow(2, exp_nk, M)

   // Compute P mod M: product_{i=0..k-1}(2^n - i) mod M
   // 2^n mod M
   pow2n := modPow(2, n_mod, M)
   var Pmod int64 = 1
   if k <= uint64(M) {
       for i := uint64(0); i < k; i++ {
           term := pow2n - int64(i%uint64(M))
           if term < 0 {
               term += M
           }
           Pmod = (Pmod * term) % M
           if Pmod == 0 {
               break
           }
       }
   } else {
       Pmod = 0
   }

   // Compute v2P and its mod phi
   // v2P = n + (k-1) - pc
   // exponent for division by 2^v2P: inv2^v2P
   v2P := int64(0)
   // v2P mod phi
   v2P = (n_mod + int64((k-1)%uint64(phi)) - int64(pc)%phi) % phi
   if v2P < 0 {
       v2P += phi
   }
   // Compute inv2^v2P mod M
   inv2 := (M + 1) / 2
   inv2v := modPow(inv2, v2P, M)

   // Compute A' mod M = (2^{n*k} - P) * inv2^v2P mod M
   A := (pow2nk - Pmod) % M
   if A < 0 {
       A += M
   }
   A = (A * inv2v) % M

   // Compute exponent for B': exp_B = (n-1)*(k-1) + pc
   e1 := (n_mod - 1) % phi
   if e1 < 0 {
       e1 += phi
   }
   e2 := int64((k-1) % uint64(phi))
   expB := (e1*e2 + int64(pc)) % phi

   B := modPow(2, expB, M)
   fmt.Println(A, B)
}
