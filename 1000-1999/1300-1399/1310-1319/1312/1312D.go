package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

// modpow computes a^e mod mod
func modpow(a, e int64) int64 {
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
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // For n<3, cannot have exactly one duplicate under mountain constraints
   if n < 3 {
       fmt.Fprintln(writer, 0)
       return
   }
   // Precompute factorials and inverse factorials up to m
   fac := make([]int64, m+1)
   ifac := make([]int64, m+1)
   fac[0] = 1
   for i := 1; i <= m; i++ {
       fac[i] = fac[i-1] * int64(i) % mod
   }
   ifac[m] = modpow(fac[m], mod-2)
   for i := m; i > 0; i-- {
       ifac[i-1] = ifac[i] * int64(i) % mod
   }
   // Compute C(m-2, n-3)
   cmb := int64(0)
   N := m - 2
   K := n - 3
   if K >= 0 && N >= K {
       cmb = fac[N] * ifac[K] % mod * ifac[N-K] % mod
   }
   // Compute m*(m-1)/2 mod mod
   // inverse of 2 modulo mod
   inv2 := int64((mod + 1) / 2)
   x := int64(m) * int64(m-1) % mod * inv2 % mod
   // 2^(n-3)
   pow2 := modpow(2, int64(n-3))
   ans := x * cmb % mod * pow2 % mod
   fmt.Fprintln(writer, ans)
}
