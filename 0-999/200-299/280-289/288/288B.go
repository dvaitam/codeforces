package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func modpow(a, e int64) int64 {
   if e == 0 {
       return 1
   }
   if e < 0 {
       // negative exponent: compute inverse
       return modpow(modpow(a, -e), MOD-2)
   }
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
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // precompute factorials up to k
   fact := make([]int64, k+1)
   invfact := make([]int64, k+1)
   fact[0] = 1
   for i := 1; i <= k; i++ {
       fact[i] = fact[i-1] * int64(i) % MOD
   }
   for i := 0; i <= k; i++ {
       invfact[i] = modpow(fact[i], MOD-2)
   }
   // sum over possible cycle sizes m
   var sum int64
   for m := 1; m <= k; m++ {
       // cycle_count = (k-1)! / (k-m)!
       cyc := fact[k-1] * invfact[k-m] % MOD
       // tree attachments for remaining k-m nodes
       var trees int64
       if m < k {
           // m * k^(k-m-1)
           trees = int64(m) * modpow(int64(k), int64(k-m-1)) % MOD
       } else {
           trees = 1
       }
       sum = (sum + cyc*trees) % MOD
   }
   // nodes > k map within themselves: (n-k)^(n-k)
   rem := modpow(int64(n-k), int64(n-k))
   ans := sum * rem % MOD
   fmt.Fprintln(writer, ans)
}
