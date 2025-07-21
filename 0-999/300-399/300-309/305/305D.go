package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MOD = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // Count jump positions S_pos = n-k-1
   Spos := n - k - 1
   jumps := make([]int, 0, m)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       d := v - u
       if d == 1 {
           continue
       } else if d == k+1 {
           jumps = append(jumps, u)
       } else {
           fmt.Println(0)
           return
       }
   }
   sort.Ints(jumps)
   P := len(jumps)
   // No possible jump positions
   if Spos <= 0 {
       if P == 0 {
           fmt.Println(1)
       } else {
           fmt.Println(0)
       }
       return
   }
   // Precompute powers of 2 up to max(Spos, k+1)
   maxN := Spos
   if k+1 > maxN {
       maxN = k + 1
   }
   pow2 := make([]int, maxN+2)
   pow2[0] = 1
   for i := 1; i <= maxN+1; i++ {
       pow2[i] = pow2[i-1] * 2 % MOD
   }
   // No initial jumps
   if P == 0 {
       if Spos <= k {
           fmt.Println(pow2[Spos])
       } else {
           // (Spos - k + 1) * 2^k
           res := int64(Spos - k + 1) * int64(pow2[k]) % MOD
           fmt.Println(res)
       }
       return
   }
   imin := jumps[0]
   imax := jumps[P-1]
   if imax-imin > k {
       fmt.Println(0)
       return
   }
   // Range of window starts L0
   low := imax - k
   if low < 1 {
       low = 1
   }
   high := imin
   res := 0
   // helper for lower_bound and upper_bound via sort.Search
   for L0 := low; L0 <= high; L0++ {
       // window [L0..R0]
       R0 := L0 + k
       if R0 > Spos {
           R0 = Spos
       }
       wlen := R0 - L0 + 1
       // count initial jumps in [L0..R0]
       lo := sort.Search(len(jumps), func(i int) bool { return jumps[i] >= L0 })
       hi := sort.Search(len(jumps), func(i int) bool { return jumps[i] > R0 }) - 1
       inCnt := 0
       if lo < len(jumps) && hi >= lo {
           inCnt = hi - lo + 1
       }
       aSz := wlen - inCnt
       var add int
       if L0 == imin {
           // L0 in initial, all aSz optional
           add = pow2[aSz]
       } else {
           // L0 not in initial, must include L0, so choose among aSz-1
           if aSz-1 >= 0 {
               add = pow2[aSz-1]
           } else {
               add = 0
           }
       }
       res = (res + add) % MOD
   }
   fmt.Println(res)
}
