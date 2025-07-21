package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ps := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       ps[i] = ps[i-1] + a[i]
   }
   // dpPrev[l][r]: best value for t-1 segments, last segment [l..r]
   const NEG_INF = -1 << 60
   dpPrev := make([][]int64, n+2)
   for i := range dpPrev {
       dpPrev[i] = make([]int64, n+2)
   }
   // t = 1: dpPrev initialized to 0 for all segments [l..r]
   // iterate for t = 2..k
   for t := 2; t <= k; t++ {
       // Build maxEndA and maxEndB for dpPrev
       maxEndA := make([]int64, n+2)
       maxEndB := make([]int64, n+2)
       for i := 0; i <= n+1; i++ {
           maxEndA[i] = NEG_INF
           maxEndB[i] = NEG_INF
       }
       for l := 1; l <= n; l++ {
           for r := l; r <= n; r++ {
               d := dpPrev[l][r]
               s := ps[r] - ps[l-1]
               if d+s > maxEndA[r] {
                   maxEndA[r] = d + s
               }
               if d-s > maxEndB[r] {
                   maxEndB[r] = d - s
               }
           }
       }
       // prefix max for bestA and bestB
       bestA := make([]int64, n+2)
       bestB := make([]int64, n+2)
       bestA[0] = NEG_INF
       bestB[0] = NEG_INF
       for i := 1; i <= n; i++ {
           bestA[i] = max(bestA[i-1], maxEndA[i])
           bestB[i] = max(bestB[i-1], maxEndB[i])
       }
       // compute dpCur for t segments
       dpCur := make([][]int64, n+2)
       for i := range dpCur {
           dpCur[i] = make([]int64, n+2)
           for j := range dpCur[i] {
               dpCur[i][j] = NEG_INF
           }
       }
       for l := 1; l <= n; l++ {
           // previous segments must end before l
           pa := bestA[l-1]
           pb := bestB[l-1]
           if pa == NEG_INF && pb == NEG_INF {
               continue
           }
           for r := l; r <= n; r++ {
               s := ps[r] - ps[l-1]
               v1 := pa - s
               v2 := pb + s
               dpCur[l][r] = max(v1, v2)
           }
       }
       dpPrev = dpCur
   }
   // answer is max over dpPrev[l][r]
   ans := int64(0)
   first := true
   for l := 1; l <= n; l++ {
       for r := l; r <= n; r++ {
           v := dpPrev[l][r]
           if first || v > ans {
               ans = v
               first = false
           }
       }
   }
   fmt.Println(ans)
}
