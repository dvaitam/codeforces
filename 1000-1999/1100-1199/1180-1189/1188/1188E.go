package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 998244353

func modpow(a, e int64) int64 {
   res := int64(1)
   base := a % MOD
   for e > 0 {
       if e&1 == 1 {
           res = res * base % MOD
       }
       base = base * base % MOD
       e >>= 1
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   a := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // Precompute factorials up to 2*k
   maxN := 2*k + 5
   fac := make([]int64, maxN)
   ifac := make([]int64, maxN)
   fac[0] = 1
   for i := 1; i < maxN; i++ {
       fac[i] = fac[i-1] * int64(i) % MOD
   }
   ifac[maxN-1] = modpow(fac[maxN-1], MOD-2)
   for i := maxN - 2; i >= 0; i-- {
       ifac[i] = ifac[i+1] * int64(i+1) % MOD
   }

   comb := func(n, r int) int64 {
       if r < 0 || r > n || n < 0 {
           return 0
       }
       return fac[n] * ifac[r] % MOD * ifac[n-r] % MOD
   }

   ans := int64(0)
   pos := 0
   // iterate r = t mod k from 0 to k-1
   for r := 0; r < k; r++ {
       for pos < k && a[pos] < r {
           pos++
       }
       if pos > r {
           continue
       }
       slack := r - pos
       // number of y_i >=0, sum y = slack is C(slack + k -1, k -1)
       ans = (ans + comb(slack+k-1, k-1)) % MOD
   }
   fmt.Fprintln(writer, ans)
}
