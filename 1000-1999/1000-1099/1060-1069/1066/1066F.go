package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type P struct { x, y int }

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   p := make([]P, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &p[i].x, &p[i].y)
   }
   // sort by max(x,y), then x asc, then y desc
   pts := p[1:]
   sort.Slice(pts, func(i, j int) bool {
       a, b := pts[i], pts[j]
       fa, fb := max(a.x, a.y), max(b.x, b.y)
       if fa != fb {
           return fa < fb
       }
       if a.x != b.x {
           return a.x < b.x
       }
       return a.y > b.y
   })
   // dp[i]: minimal cost reaching block endpoint i
   const INF = int64(1e18)
   dp := make([]int64, n+2)
   for i := range dp {
       dp[i] = INF
   }
   dp[0] = 0
   pl, pr := 0, 0
   // helper functions using p slice closure
   le := func(d int) int {
       return max(p[d].x, p[d].y)
   }
   get := func(a, b int) int {
       return abs(p[a].x-p[b].x) + abs(p[a].y-p[b].y)
   }
   for pr < n {
       d := pr + 1
       for d < n && le(d+1) == le(pr+1) {
           d++
       }
       // compute dp at segment end and start
       c1 := dp[pl] + int64(get(pl, pr+1))
       c2 := dp[pr] + int64(get(pr, pr+1))
       costSeg := int64(get(pr+1, d))
       if c1 < c2 {
           dp[d] = c1 + costSeg
       } else {
           dp[d] = c2 + costSeg
       }
       c3 := dp[pl] + int64(get(pl, d))
       c4 := dp[pr] + int64(get(pr, d))
       if c3 < c4 {
           dp[pr+1] = c3 + costSeg
       } else {
           dp[pr+1] = c4 + costSeg
       }
       pl = pr + 1
       pr = d
   }
   // result is min(dp[pl], dp[pr])
   res := dp[pl]
   if dp[pr] < res {
       res = dp[pr]
   }
   fmt.Println(res)
}
