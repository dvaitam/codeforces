package main

import (
   "bufio"
   "fmt"
   "os"
)

type DSU struct {
   p, r []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
       r[i] = 0
   }
   return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return
   }
   if d.r[rx] < d.r[ry] {
       d.p[rx] = ry
   } else if d.r[ry] < d.r[rx] {
       d.p[ry] = rx
   } else {
       d.p[ry] = rx
       d.r[rx]++
   }
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, W int
   fmt.Fscan(reader, &n, &m, &W)
   w := make([]int, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   dsu := NewDSU(n)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       // convert to 0-based
       dsu.Union(x-1, y-1)
   }
   // group members by root
   groups := make(map[int][]int)
   for i := 0; i < n; i++ {
       r := dsu.Find(i)
       groups[r] = append(groups[r], i)
   }
   // dp[j] = max beauty with weight j
   const inf = int64(1e18)
   dp := make([]int64, W+1)
   for i := 1; i <= W; i++ {
       dp[i] = -inf
   }
   // process each group
   for _, members := range groups {
       // compute group total
       var sumW int
       var sumB int64
       for _, i := range members {
           sumW += w[i]
           sumB += b[i]
       }
       // copy dp
       dp2 := make([]int64, W+1)
       for j := 0; j <= W; j++ {
           dp2[j] = dp[j]
       }
       // option: take whole group
       if sumW <= W {
           for j := sumW; j <= W; j++ {
               if dp[j-sumW] > -inf {
                   dp2[j] = max(dp2[j], dp[j-sumW]+sumB)
               }
           }
       }
       // option: take one member
       for _, i := range members {
           wi := w[i]
           bi := b[i]
           if wi > W {
               continue
           }
           for j := wi; j <= W; j++ {
               if dp[j-wi] > -inf {
                   dp2[j] = max(dp2[j], dp[j-wi]+bi)
               }
           }
       }
       dp = dp2
   }
   // find max
   var ans int64
   for j := 0; j <= W; j++ {
       if dp[j] > ans {
           ans = dp[j]
       }
   }
   fmt.Fprintln(writer, ans)
}
