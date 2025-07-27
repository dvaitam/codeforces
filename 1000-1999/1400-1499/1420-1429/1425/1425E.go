package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, K int
   fmt.Fscan(in, &N, &K)
   A := make([]int64, N+1)
   D := make([]int64, N+1)
   for i := 1; i <= N; i++ {
       fmt.Fscan(in, &A[i])
   }
   for i := 1; i <= N; i++ {
       fmt.Fscan(in, &D[i])
   }
   // prefix sums
   ps := make([]int64, N+1)
   for i := 1; i <= N; i++ {
       ps[i] = ps[i-1] + A[i]
   }
   // helper to run DP with penalty C per segment
   run := func(C int64) (cnt int, dpVal int64) {
       var dpPrev int64 = 0
       var cntPrev int = 0
       // best M = max(dp[j-1] - ps[j-1] - D[j])
       var bestM int64 = -1<<60
       var bestCnt int = 0
       for r := 1; r <= N; r++ {
           // update M for segment starting at r
           m := dpPrev - ps[r-1] - D[r]
           if m > bestM || (m == bestM && cntPrev < bestCnt) {
               bestM = m
               bestCnt = cntPrev
           }
           // consider ending segment at r
           raw := ps[r] + bestM
           var dpCur int64
           var cntCur int
           cand := raw - C
           if cand > dpPrev {
               dpCur = cand
               cntCur = bestCnt + 1
           } else if cand < dpPrev {
               dpCur = dpPrev
               cntCur = cntPrev
           } else {
               dpCur = dpPrev
               // tie: choose fewer segments
               if cntPrev < bestCnt+1 {
                   cntCur = cntPrev
               } else {
                   cntCur = bestCnt + 1
               }
           }
           dpPrev = dpCur
           cntPrev = cntCur
       }
       return cntPrev, dpPrev
   }
   // binary search penalty C >=0
   var l, r int64 = 0, ps[N]
   for l < r {
       mid := (l + r + 1) >> 1
       cnt, _ := run(mid)
       if cnt >= K+1 {
           l = mid
       } else {
           r = mid - 1
       }
   }
   C := l
   _, dpv := run(C)
   // result is dpv + C*(K+1)
   res := dpv + C*int64(K+1)
   fmt.Println(res)
}
