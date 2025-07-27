package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 998244353

func modPow(a, e int64) int64 {
   res := int64(1)
   a %= MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * a % MOD
       }
       a = a * a % MOD
       e >>= 1
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, 2*n)
   for i := 0; i < 2*n; i++ {
       fmt.Fscan(in, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   // sum of last n minus sum of first n
   var s1, s2 int64
   for i := 0; i < n; i++ {
       s1 = (s1 + a[i]) % MOD
   }
   for i := n; i < 2*n; i++ {
       s2 = (s2 + a[i]) % MOD
   }
   diff := (s2 - s1) % MOD
   if diff < 0 {
       diff += MOD
   }
   // compute C(2n, n)
   m := 2 * n
   fact := make([]int64, m+1)
   inv := make([]int64, m+1)
   fact[0] = 1
   for i := 1; i <= m; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   inv[m] = modPow(fact[m], MOD-2)
   for i := m; i > 0; i-- {
       inv[i-1] = inv[i] * int64(i) % MOD
   }
   comb := fact[m] * inv[n] % MOD * inv[n] % MOD
   ans := diff * comb % MOD
   fmt.Fprintln(out, ans)
}
