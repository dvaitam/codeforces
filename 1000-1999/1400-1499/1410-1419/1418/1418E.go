package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 998244353

func modInv(n int, inv []int) int {
   return inv[n]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)
   d := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &d[i])
   }
   sort.Ints(d)
   // prefix sums mod MOD
   ps := make([]int, n+1)
   for i := 0; i < n; i++ {
       ps[i+1] = (ps[i] + d[i]) % MOD
   }
   // precompute inverses up to n+1
   inv := make([]int, n+2)
   inv[1] = 1
   for i := 2; i <= n+1; i++ {
       inv[i] = MOD - int((MOD/int64(i))*int64(inv[MOD%i])%MOD)
   }

   for qi := 0; qi < m; qi++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       // count strong monsters: d[i] >= b
       idx := sort.Search(n, func(i int) bool { return d[i] >= b })
       K := n - idx
       if K < a {
           fmt.Fprintln(out, 0)
           continue
       }
       sumWeak := ps[idx]
       sumStrong := (ps[n] - ps[idx] + MOD) % MOD
       // term1: strong
       // P_s = (K-a)/K
       term1 := int((int64(K-a) * int64(modInv(K, inv)) % MOD) * int64(sumStrong) % MOD)
       // term2: weak
       // P_w = (K+1-a)/(K+1)
       term2 := int((int64(K+1-a) * int64(modInv(K+1, inv)) % MOD) * int64(sumWeak) % MOD)
       ans := term1 + term2
       if ans >= MOD {
           ans -= MOD
       }
       fmt.Fprintln(out, ans)
   }
}
