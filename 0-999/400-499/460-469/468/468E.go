package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int
   fmt.Fscan(reader, &n, &k)
   xs := make([]int, k)
   ys := make([]int, k)
   ds := make([]int64, k)
   for i := 0; i < k; i++ {
       var x, y int
       var w int64
       fmt.Fscan(reader, &x, &y, &w)
       xs[i] = x
       ys[i] = y
       ds[i] = (w - 1 + mod) % mod
   }
   // Precompute conflict masks: edges sharing row or column
   conflictMask := make([]uint64, k)
   for i := 0; i < k; i++ {
       var m uint64
       for j := 0; j < k; j++ {
           if xs[i] == xs[j] || ys[i] == ys[j] {
               m |= 1 << uint(j)
           }
       }
       conflictMask[i] = m
   }
   // Memoization for matching generating function
   memo := make(map[uint64][]int64)
   base := make([]int64, k+1)
   base[0] = 1
   memo[0] = base

   var solve func(mask uint64) []int64
   solve = func(mask uint64) []int64 {
       if res, ok := memo[mask]; ok {
           return res
       }
       // pick an edge
       e := uint(bits.TrailingZeros64(mask))
       // exclude edge e
       maskWithout := mask & (mask - 1)
       res0 := solve(maskWithout)
       // include edge e: remove conflicting edges
       maskInclude := mask &^ conflictMask[e]
       res1 := solve(maskInclude)
       // combine results
       res := make([]int64, k+1)
       res[0] = res0[0]
       for t := 1; t <= k; t++ {
           v := res0[t]
           v += ds[e] * res1[t-1] % mod
           if v >= mod {
               v -= mod
           }
           res[t] = v
       }
       memo[mask] = res
       return res
   }

   // full set of edges
   var fullMask uint64
   for i := 0; i < k; i++ {
       fullMask |= 1 << uint(i)
   }
   eSum := solve(fullMask)

   // factorials up to n
   fact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   var ans int64
   for t, v := range eSum {
       if v == 0 || t > n {
           continue
       }
       ans = (ans + v*fact[n-t]) % mod
   }
   fmt.Println(ans)
}
