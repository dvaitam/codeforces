package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func modPow(x, y int64) int64 {
   res := int64(1)
   x %= mod
   for y > 0 {
       if y&1 == 1 {
           res = res * x % mod
       }
       x = x * x % mod
       y >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   // edge cases handled by general formula
   if m == 1 {
       if k == 0 {
           fmt.Fprintln(writer, 1)
       } else {
           fmt.Fprintln(writer, 0)
       }
       return
   }
   if k < 0 || k > n-1 {
       fmt.Fprintln(writer, 0)
       return
   }
   // precompute factorials and inverse factorials up to n-1
   maxN := n - 1
   fact := make([]int64, maxN+1)
   invFact := make([]int64, maxN+1)
   fact[0] = 1
   for i := 1; i <= maxN; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   invFact[maxN] = modPow(fact[maxN], mod-2)
   for i := maxN; i > 0; i-- {
       invFact[i-1] = invFact[i] * int64(i) % mod
   }
   // compute C(n-1, k)
   comb := fact[maxN] * invFact[k] % mod * invFact[maxN-k] % mod
   // m * C(n-1,k) * (m-1)^k
   ans := int64(m) * comb % mod * modPow(int64(m-1), int64(k)) % mod
   fmt.Fprintln(writer, ans)
}
