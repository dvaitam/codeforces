package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &c[i])
   }
   p := make([][]int64, n)
   for i := 0; i < n; i++ {
       p[i] = make([]int64, m+1)
       for j := 1; j <= m; j++ {
           fmt.Fscan(in, &p[i][j])
       }
   }
   const inf = int64(1e18)
   // dpPrev[j][l]: min cost for first i-1 trees, beauty j, ending color l
   dpPrev := make([][]int64, k+1)
   dpCur := make([][]int64, k+1)
   for j := 0; j <= k; j++ {
       dpPrev[j] = make([]int64, m+1)
       dpCur[j] = make([]int64, m+1)
       for l := 0; l <= m; l++ {
           dpPrev[j][l] = inf
           dpCur[j][l] = inf
       }
   }
   // before any trees, zero segments cost 0 for any last color
   for l := 1; l <= m; l++ {
       dpPrev[0][l] = 0
   }
   // process each tree
   for i := 0; i < n; i++ {
       // reset dpCur
       for j := 0; j <= k; j++ {
           for l := 1; l <= m; l++ {
               dpCur[j][l] = inf
           }
       }
       // precompute best and second best for dpPrev
       best1 := make([]int64, k+1)
       best2 := make([]int64, k+1)
       best1Color := make([]int, k+1)
       for j := 0; j < k; j++ {
           b1, b2 := inf, inf
           bc := 0
           for l := 1; l <= m; l++ {
               v := dpPrev[j][l]
               if v < b1 {
                   b2 = b1
                   b1 = v
                   bc = l
               } else if v < b2 {
                   b2 = v
               }
           }
           best1[j] = b1
           best2[j] = b2
           best1Color[j] = bc
       }
       // for each desired segment count
       for j := 1; j <= k; j++ {
           for l := 1; l <= m; l++ {
               // cost to paint or keep
               cost := inf
               if c[i] == 0 {
                   cost = p[i][l]
               } else if c[i] == l {
                   cost = 0
               } else {
                   continue
               }
               // continuation of segment
               v := dpPrev[j][l]
               // start new segment
               if j > 0 {
                   // use best of j-1 segments
                   b := best1[j-1]
                   if best1Color[j-1] == l {
                       b = best2[j-1]
                   }
                   if b < v {
                       v = b
                   }
               }
               if v < inf {
                   dpCur[j][l] = v + cost
               }
           }
       }
       // swap dpPrev and dpCur
       dpPrev, dpCur = dpCur, dpPrev
   }
   // answer is min over dpPrev[k][l]
   ans := inf
   for l := 1; l <= m; l++ {
       if dpPrev[k][l] < ans {
           ans = dpPrev[k][l]
       }
   }
   if ans >= inf/2 {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
