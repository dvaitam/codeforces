package main

import (
   "bufio"
   "fmt"
   "os"
)

// modInv computes modular inverse of a under mod
func modInv(a, mod int64) int64 {
   var b = mod - 2
   res := int64(1)
   for b > 0 {
       if b&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       b >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var T int64
   fmt.Fscan(reader, &n, &T)
   t := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i])
   }
   mod := int64(1e9 + 7)
   // precompute factorials and inverse factorials up to n
   fac := make([]int64, n+1)
   ifac := make([]int64, n+1)
   fac[0] = 1
   for i := 1; i <= n; i++ {
       fac[i] = fac[i-1] * int64(i) % mod
   }
   ifac[n] = modInv(fac[n], mod)
   for i := n; i > 0; i-- {
       ifac[i-1] = ifac[i] * int64(i) % mod
   }
   // precompute inverse powers of 2
   inv2 := modInv(2, mod)
   invPow2 := make([]int64, n+1)
   invPow2[0] = 1
   for i := 1; i <= n; i++ {
       invPow2[i] = invPow2[i-1] * inv2 % mod
   }
   var res int64 = 0
   var sumT int64 = 0
   for k := 1; k <= n; k++ {
       sumT += t[k-1]
       D := T - sumT
       if D < 0 {
           break
       }
       if D >= int64(k) {
           res = (res + 1) % mod
       } else {
           var s int64 = 0
           lim := int(D)
           for i := 0; i <= lim; i++ {
               // C(k, i) = fac[k] * ifac[i] * ifac[k-i]
               s = (s + fac[k] * ifac[i] % mod * ifac[k-i] % mod) % mod
           }
           res = (res + s * invPow2[k] % mod) % mod
       }
   }
   fmt.Println(res)
}
