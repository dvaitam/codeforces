package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000009
const MAXN = 4005

var fac [MAXN]int64
var invfac [MAXN]int64

// modular exponentiation
func modpow(a, e int64) int64 {
   res := int64(1)
   for e > 0 {
       if e&1 == 1 {
           res = res * a % mod
       }
       a = a * a % mod
       e >>= 1
   }
   return res
}

// combination nCk
func comb(n, k int) int64 {
   if k < 0 || k > n {
       return 0
   }
   return fac[n] * (invfac[k] * invfac[n-k] % mod) % mod
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, w, b int
   if _, err := fmt.Fscan(reader, &n, &w, &b); err != nil {
       return
   }
   // precompute factorials and inverse factorials
   fac[0] = 1
   for i := 1; i < MAXN; i++ {
       fac[i] = fac[i-1] * int64(i) % mod
   }
   // invfac[MAXN-1] via Fermat
   invfac[MAXN-1] = modpow(fac[MAXN-1], mod-2)
   for i := MAXN - 1; i > 0; i-- {
       invfac[i-1] = invfac[i] * int64(i) % mod
   }

   // precompute w! and b!
   wf := fac[w]
   bf := fac[b]

   var sum int64
   // t = total white days = x+z, t from 2 to n-1
   for t := 2; t <= n-1; t++ {
       y := n - t // black days
       if y < 1 {
           continue
       }
       // number of ways to choose x,z positive with x+z=t is (t-1)
       waysStripes := int64(t - 1)
       c1 := comb(w-1, t-1)
       c2 := comb(b-1, y-1)
       if c1 == 0 || c2 == 0 {
           continue
       }
       sum = (sum + waysStripes * c1 % mod * c2) % mod
   }
   ans := wf * bf % mod * sum % mod
   fmt.Fprint(writer, ans)
}
