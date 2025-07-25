package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   const INF = int64(4e18)
   dpPrev := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dpPrev[i] = INF
   }
   // dpPrev[0] = 0

   type seg struct {
       max   int
       left  int
       best  int64
   }

   // Build RMQ (sparse table) for range maximum queries
   log := make([]int, n+2)
   log[1] = 0
   for i := 2; i <= n+1; i++ {
       log[i] = log[i/2] + 1
   }
   K := log[n] + 1
   st := make([][]int, K)
   st[0] = make([]int, n+1)
   for i := 1; i <= n; i++ {
       st[0][i] = a[i]
   }
   for j := 1; j < K; j++ {
       st[j] = make([]int, n+1)
       shift := 1 << (j - 1)
       for i := 1; i+shift <= n; i++ {
           x := st[j-1][i]
           y := st[j-1][i+shift]
           if x > y {
               st[j][i] = x
           } else {
               st[j][i] = y
           }
       }
   }
   // helper to get max in [l..r]
   rmq := func(l, r int) int {
       if l > r {
           return 0
       }
       length := r - l + 1
       j := log[length]
       x := st[j][l]
       y := st[j][r-(1<<j)+1]
       if x > y {
           return x
       }
       return y
   }
   dpCur := make([]int64, n+1)
   var solve func(int, int, int, int)
   solve = func(l, r, optL, optR int) {
       if l > r {
           return
       }
       mid := (l + r) >> 1
       bestP := -1
       dpCur[mid] = INF
       // p from optL to min(mid-1, optR)
       end := mid - 1
       if end > optR {
           end = optR
       }
       for p := optL; p <= end; p++ {
           // cost of segment (p+1..mid)
           m := rmq(p+1, mid)
           cost := dpPrev[p] + int64(mid-p)*int64(m)
           if cost < dpCur[mid] {
               dpCur[mid] = cost
               bestP = p
           }
       }
       if bestP < 0 {
           bestP = optL
       }
       solve(l, mid-1, optL, bestP)
       solve(mid+1, r, bestP, optR)
   }
   // DP layers with divide and conquer optimization
   for segCnt := 1; segCnt <= k; segCnt++ {
       // positions < segCnt cannot have segCnt segments
       for i := 0; i < segCnt; i++ {
           dpCur[i] = INF
       }
       // compute dpCur[segCnt..n]
       solve(segCnt, n, segCnt-1, n-1)
       // swap dpPrev and dpCur
       dpPrev, dpCur = dpCur, dpPrev
   }
   // result
   fmt.Fprintln(writer, dpPrev[n])
}
