package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

const mod = 998244353

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   k := n - 1
   L := make([]uint64, k)
   D := make([]uint64, k)
   var xorA uint64
   for i := 0; i < n; i++ {
       xorA ^= a[i]
       if i < n-1 {
           L[i] = a[i]
           D[i] = a[i+1] - a[i]
       }
   }
   // Transform xor target
   var xorL uint64
   for i := 0; i < k; i++ {
       xorL ^= L[i]
   }
   target := xorA ^ xorL

   // Precompute Bmask for each bit
   const maxB = 60
   Bmask := make([]int, maxB)
   for b := 0; b < maxB; b++ {
       var m int
       for j := 0; j < k; j++ {
           if (D[j]>>uint(b))&1 != 0 {
               m |= 1 << j
           }
       }
       Bmask[b] = m
   }
   // Precompute powers of 2 up to k
   pow2 := make([]int, k+1)
   pow2[0] = 1
   for i := 1; i <= k; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   // DP over masks
   S := 1 << k
   dp := make([]int, S)
   full := S - 1
   dp[full] = 1
   // process bits from high to low
   for b := maxB - 1; b >= 0; b-- {
       t := int((target >> uint(b)) & 1)
       bm := Bmask[b]
       dp2 := make([]int, S)
       for mask := 0; mask < S; mask++ {
           v := dp[mask]
           if v == 0 {
               continue
           }
           Bcap := mask & bm
           Zmask := mask &^ bm
           u := k - bits.OnesCount(uint(mask))
           // iterate subsets of Bcap
           if Bcap == 0 {
               // only s=0
               var cnt int
               if u > 0 {
                   cnt = pow2[u-1]
               } else {
                   if t == 0 {
                       cnt = 1
                   } else {
                       cnt = 0
                   }
               }
               mask2 := Zmask
               dp2[mask2] = (dp2[mask2] + v*cnt) % mod
           } else {
               // general
               // precompute cntU if u>0
               cntU := 0
               if u > 0 {
                   cntU = pow2[u-1]
               }
               for s := Bcap; ; s = (s - 1) & Bcap {
                   p := bits.OnesCount(uint(s)) & 1
                   var cnt int
                   if u > 0 {
                       cnt = cntU
                   } else {
                       if (t ^ p) == 0 {
                           cnt = 1
                       } else {
                           cnt = 0
                       }
                   }
                   if u > 0 || (t ^ p) == 0 {
                       mask2 := Zmask | s
                       dp2[mask2] = (dp2[mask2] + v*cnt) % mod
                   }
                   if s == 0 {
                       break
                   }
               }
           }
       }
       dp = dp2
   }
   // sum all dp
   var ans int
   for _, v := range dp {
       ans = (ans + v) % mod
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
