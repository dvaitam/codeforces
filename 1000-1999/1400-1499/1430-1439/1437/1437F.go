package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func modinv(a, mod int) int {
   b := mod
   u, v := 1, 0
   for b != 0 {
       t := a / b
       a, b = b, a - t*b
       u, v = v, u - t*v
   }
   if u < 0 {
       u += mod
   }
   return u
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   mod := 998244353
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // precompute factorials and invfacts
   fact := make([]int, n+1)
   invf := make([]int, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = int(int64(fact[i-1]) * int64(i) % int64(mod))
   }
   invf[n] = modinv(fact[n], mod)
   for i := n; i > 0; i-- {
       invf[i-1] = int(int64(invf[i]) * int64(i) % int64(mod))
   }
   ans := 0
   // consider suffix starting at f
   for f := 1; f <= n; f++ {
       ok := true
       for i := f; i < n; i++ {
           if a[i] < 2*a[i-1] {
               ok = false
               break
           }
       }
       if !ok {
           continue
       }
       m := n - f + 1
       F := f - 1
       c := make([]int, m+1)
       // compute release times
       for k := 0; k < f-1; k++ {
           // binary search in a[f-1..n-1]
           lo, hi := f-1, n-1
           pos := -1
           for lo <= hi {
               mid := (lo + hi) / 2
               if a[mid] >= 2*a[k] {
                   pos = mid
                   hi = mid - 1
               } else {
                   lo = mid + 1
               }
           }
           if pos < 0 {
               ok = false
               break
           }
           t := pos - (f - 1) + 1
           if t < 1 || t > m {
               ok = false
               break
           }
           c[t]++
       }
       if !ok {
           continue
       }
       dp := make([]int, F+1)
       dp[0] = 1
       avail := 0
       for j := 1; j <= m; j++ {
           avail += c[j]
           ndp := make([]int, F+1)
           for used := 0; used <= F; used++ {
               v := dp[used]
               if v == 0 {
                   continue
               }
               maxk := avail - used
               if maxk > F-used {
                   maxk = F - used
               }
               for k := 0; k <= maxk; k++ {
                   // P(avail-used, k) = fact[avail-used] / fact[avail-used-k]
                   ways := int(int64(fact[avail-used]) * int64(invf[avail-used-k]) % int64(mod))
                   ndp[used+k] = (ndp[used+k] + v*ways) % mod
               }
           }
           dp = ndp
       }
       ans = (ans + dp[F]) % mod
   }
   fmt.Println(ans)
}
