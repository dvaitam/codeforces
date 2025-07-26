package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF64 = uint64(2e18)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   var k uint64
   fmt.Fscan(in, &n, &m, &k)
   var s string
   fmt.Fscan(in, &s)
   if m == 1 {
       fmt.Println(s)
       return
   }
   // prepare rolling hash
   base := uint64(1315423911)
   n0 := len(s)
   h := make([]uint64, n0+1)
   pows := make([]uint64, n0+1)
   pows[0] = 1
   for i := 0; i < n0; i++ {
       h[i+1] = h[i]*base + uint64(s[i])
       pows[i+1] = pows[i] * base
   }
   // function to get hash of s[l:r] inclusive, 0-indexed
   hashSeg := func(l, r int) uint64 {
       return h[r+1] - h[l]*pows[r-l+1]
   }
   // cur prefix
   cur := ""
   // DP buffers
   dpPrev := make([]uint64, n0+1)
   dpCurr := make([]uint64, n0+1)
   // lim and events
   lim := make([]int, n0)
   addList := make([][]int, n0+1)

   // function to compute f(t) = number of splits where minimal substring >= t
   var fcount func(t string) uint64
   fcount = func(t string) uint64 {
       L := len(t)
       // compute hash of t
       ht := make([]uint64, L+1)
       ht[0] = 0
       for i := 0; i < L; i++ {
           ht[i+1] = ht[i]*base + uint64(t[i])
       }
       // LCP function between s[l:] and t
       lcp := func(l int) int {
           // binary search on [0, min]
           lo, hi := 0, n0-l
           if hi > L {
               hi = L
           }
           for lo < hi {
               mid := (lo + hi + 1) >> 1
               if hashSeg(l, l+mid-1) == ht[mid] {
                   lo = mid
               } else {
                   hi = mid - 1
               }
           }
           return lo
       }
       // compute lim array
       for i := range addList {
           addList[i] = addList[i][:0]
       }
       for i := 0; i < n0; i++ {
           Lcp := lcp(i)
           if Lcp < L {
               // suffix ended or mismatch
               if i+Lcp >= n0 {
                   lim[i] = n0 + 1
               } else if s[i+Lcp] > t[Lcp] {
                   lim[i] = i + Lcp
               } else {
                   lim[i] = n0 + 1
               }
           } else {
               // full match of t
               lim[i] = i + L - 1
           }
           if lim[i] < n0 {
               addList[lim[i]] = append(addList[lim[i]], i+1) // 1-based l
           }
       }
       // init dp
       for i := 0; i <= n0; i++ {
           dpPrev[i] = 0
       }
       dpPrev[0] = 1
       // run DP for j=1..m
       for j := 1; j <= m; j++ {
           // clear curr
           for i := 0; i <= n0; i++ {
               dpCurr[i] = 0
           }
           var curSum uint64
           // for i from 1..n0
           for i := 1; i <= n0; i++ {
               // add events
               for _, l := range addList[i-1] {
                   // segment from l to i-1, dpPrev at l-1
                   curSum += dpPrev[l-1]
                   if curSum > k {
                       curSum = k + 1
                   }
               }
               if j == 1 {
                   // only one segment: dpCurr[i] = curSum
                   dpCurr[i] = curSum
               } else {
                   dpCurr[i] = curSum
               }
           }
           // swap
           dpPrev, dpCurr = dpCurr, dpPrev
       }
       res := dpPrev[n0]
       if res > k {
           res = k + 1
       }
       return res
   }

   // build answer greedily
   for {
       found := false
       for c := byte('z'); c >= 'a'; c-- {
           cand := cur + string(c)
           if fcount(cand) >= k {
               cur = cand
               found = true
               break
           }
           if c == 'a' {
               break
           }
       }
       if !found {
           break
       }
   }
   fmt.Println(cur)
}
