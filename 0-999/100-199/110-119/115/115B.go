package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
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
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &grid[i])
   }
   // find rows with weeds and their min/max columns
   type seg struct{ r, l, rr int }
   weeds := []seg{}
   for i := 0; i < n; i++ {
       l, r := m+1, 0
       for j := 0; j < m; j++ {
           if grid[i][j] == 'W' {
               if j+1 < l {
                   l = j + 1
               }
               if j+1 > r {
                   r = j + 1
               }
           }
       }
       if r > 0 {
           weeds = append(weeds, seg{i + 1, l, r})
       }
   }
   if len(weeds) == 0 {
       fmt.Println(0)
       return
   }
   k := len(weeds)
   // dp[t][0] = end at l, dp[t][1] = end at r
   const INF = 1e18
   dp0 := make([]int, k)
   dp1 := make([]int, k)
   // initial for first weed row
   f := weeds[0]
   // vertical moves down to first weed row
   base := f.r - 1
   // horizontal cost to end at L and R
   // end at L
   dp0[0] = base + (2*f.rr - f.l - 1)
   // end at R
   dp1[0] = base + (f.rr - 1)
   // DP over subsequent weed rows
   for t := 1; t < k; t++ {
       c := weeds[t]
       p := weeds[t-1]
       gap := c.r - p.r
       nd0, nd1 := int(INF), int(INF)
       // from previous end at L (p.l)
       start := p.l
       // to current L
       costL := min(abs(start-c.l)+2*(c.rr-c.l), abs(start-c.rr)+(c.rr-c.l))
       // to current R
       costR := min(abs(start-c.rr)+2*(c.rr-c.l), abs(start-c.l)+(c.rr-c.l))
       nd0 = min(nd0, dp0[t-1] + gap + costL)
       nd1 = min(nd1, dp0[t-1] + gap + costR)
       // from previous end at R (p.rr)
       start = p.rr
       costL = min(abs(start-c.l)+2*(c.rr-c.l), abs(start-c.rr)+(c.rr-c.l))
       costR = min(abs(start-c.rr)+2*(c.rr-c.l), abs(start-c.l)+(c.rr-c.l))
       nd0 = min(nd0, dp1[t-1] + gap + costL)
       nd1 = min(nd1, dp1[t-1] + gap + costR)
       dp0[t], dp1[t] = nd0, nd1
   }
   res := dp0[k-1]
   if dp1[k-1] < res {
       res = dp1[k-1]
   }
   fmt.Println(res)
}
