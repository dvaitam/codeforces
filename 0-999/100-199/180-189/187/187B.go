package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = int64(4e18)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, r int
   if _, err := fmt.Fscan(in, &n, &m, &r); err != nil {
       return
   }
   // Read car distance matrices and floyd-warshall per car
   // distC[c][i][j]
   distC := make([][][]int64, m)
   for c := 0; c < m; c++ {
       dist := make([][]int64, n)
       for i := 0; i < n; i++ {
           dist[i] = make([]int64, n)
           for j := 0; j < n; j++ {
               var x int64
               fmt.Fscan(in, &x)
               dist[i][j] = x
           }
       }
       // floyd
       for k := 0; k < n; k++ {
           dk := dist[k]
           for i := 0; i < n; i++ {
               di := dist[i]
               ik := di[k]
               for j := 0; j < n; j++ {
                   // di[j] = min(di[j], di[k] + dk[j])
                   if v := ik + dk[j]; v < di[j] {
                       di[j] = v
                   }
               }
           }
       }
       distC[c] = dist
   }
   // Build D0 matrix: best single-segment times
   D0 := make([][]int64, n)
   for i := 0; i < n; i++ {
       D0[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           best := INF
           for c := 0; c < m; c++ {
               if v := distC[c][i][j]; v < best {
                   best = v
               }
           }
           D0[i][j] = best
       }
   }
   // Read queries
   type Q struct{ idx, t, seg int }
   groups := make(map[int][]Q)
   answers := make([]int64, r)
   for i := 0; i < r; i++ {
       var s, t, k int
       fmt.Fscan(in, &s, &t, &k)
       s--;
       t--;
       seg := k + 1
       groups[s] = append(groups[s], Q{i, t, seg})
   }
   // Process per start city
   for s, qs := range groups {
       // find max segments
       maxSeg := 0
       for _, q := range qs {
           if q.seg > maxSeg {
               maxSeg = q.seg
           }
       }
       // bucket queries by seg
       bySeg := make([][]Q, maxSeg+1)
       for _, q := range qs {
           if q.seg <= maxSeg {
               bySeg[q.seg] = append(bySeg[q.seg], q)
           }
       }
       // dpPrev = dp[seg-1], dpCurr = dp[seg]
       dpPrev := make([]int64, n)
       for i := 0; i < n; i++ {
           dpPrev[i] = INF
       }
       dpPrev[s] = 0
       // iterate seg from 1 to maxSeg
       dpCurr := make([]int64, n)
       for seg := 1; seg <= maxSeg; seg++ {
           // compute dpCurr[j] = min over i dpPrev[i] + D0[i][j]
           for j := 0; j < n; j++ {
               dpCurr[j] = INF
           }
           for i := 0; i < n; i++ {
               pi := dpPrev[i]
               if pi == INF {
                   continue
               }
               di := D0[i]
               for j := 0; j < n; j++ {
                   if v := pi + di[j]; v < dpCurr[j] {
                       dpCurr[j] = v
                   }
               }
           }
           // answer queries with this seg
           for _, q := range bySeg[seg] {
               answers[q.idx] = dpCurr[q.t]
           }
           // swap dpPrev, dpCurr
           dpPrev, dpCurr = dpCurr, dpPrev
       }
   }
   // output answers
   for i := 0; i < r; i++ {
       fmt.Fprintln(out, answers[i])
   }
}
