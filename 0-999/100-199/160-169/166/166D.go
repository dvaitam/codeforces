package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   Inf int64 = 100000000000
   B   int64 = 200000
)

type item struct {
   val int
   idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   c := make([]int, n)
   ax0 := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i], &ax0[i])
   }
   // compute second set of ax values
   ax1 := make([]int, n)
   for i := 0; i < n; i++ {
       ax1[i] = ax0[i] - 1
   }
   var m int
   fmt.Fscan(reader, &m)
   d := make([]int, m)
   l := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &d[i], &l[i])
   }
   // build combined ax for compression
   N := 2*n + m
   ax := make([]int, N)
   for i := 0; i < n; i++ {
       ax[i] = ax0[i]
       ax[i+n] = ax1[i]
   }
   for i := 0; i < m; i++ {
       ax[2*n+i] = l[i]
   }
   // coordinate compression
   items := make([]item, N)
   for i := 0; i < N; i++ {
       items[i] = item{val: ax[i], idx: i}
   }
   sort.Slice(items, func(i, j int) bool { return items[i].val < items[j].val })
   comp := make([]int, N)
   M := 1
   comp[items[0].idx] = 1
   for i := 1; i < N; i++ {
       if items[i].val != items[i-1].val {
           M++
       }
       comp[items[i].idx] = M
   }
   // prepare helpers
   st := make([]int, M+2)
   for i := range st {
       st[i] = -1
   }
   for i := 0; i < n; i++ {
       st[comp[i+n]] = i
   }
   mx := make([][2]int, M+2)
   for i := range mx {
       mx[i][0], mx[i][1] = -1, -1
   }
   // find top two d for each compressed coord of l
   for i := 0; i < m; i++ {
       j := comp[2*n+i]
       if mx[j][0] == -1 || d[i] > d[mx[j][0]] {
           mx[j][1] = mx[j][0]
           mx[j][0] = i
       } else if mx[j][1] == -1 || d[i] > d[mx[j][1]] {
           mx[j][1] = i
       }
   }
   // dp arrays
   dp := make([][2]int64, M+2)
   pre := make([][2]int, M+2)
   sel := make([][2]int64, M+2)
   for i := 0; i < M+2; i++ {
       dp[i][0], dp[i][1] = -Inf, -Inf
       sel[i][0], sel[i][1] = -1, -1
   }
   dp[1][0] = 0
   // DP
   for i := 1; i <= M; i++ {
       for j := 0; j < 2; j++ {
           if dp[i][j] == -Inf {
               continue
           }
           // skip selection
           if dp[i+1][0] < dp[i][j] {
               dp[i+1][0] = dp[i][j]
               pre[i+1][0] = j
               sel[i+1][0] = -1
           }
           // take with state j
           // stay in state 0
           if j == 0 {
               if mx[i][0] != -1 && d[mx[i][0]] >= c[st[i]] {
                   val := dp[i][j] + int64(c[st[i]])
                   if dp[i+1][0] < val {
                       dp[i+1][0] = val
                       pre[i+1][0] = j
                       sel[i+1][0] = int64(st[i])*B + int64(mx[i][0])
                   }
               }
           } else {
               // j == 1
               xx := -1
               if mx[i][1] != -1 && d[mx[i][1]] >= c[st[i-1]] {
                   xx = mx[i][0]
               } else {
                   xx = mx[i][1]
               }
               if xx != -1 && d[xx] >= c[st[i]] {
                   val := dp[i][j] + int64(c[st[i]])
                   if dp[i+1][0] < val {
                       dp[i+1][0] = val
                       pre[i+1][0] = j
                       sel[i+1][0] = int64(st[i])*B + int64(xx)
                   }
               }
           }
           // move to state 1
           if mx[i+1][0] != -1 && d[mx[i+1][0]] >= c[st[i]] {
               // choose best for state 1
               choose := mx[i+1][0]
               if mx[i+1][1] != -1 && d[mx[i+1][1]] >= c[st[i]] {
                   choose = mx[i+1][1]
               }
               val := dp[i][j] + int64(c[st[i]])
               if dp[i+1][1] < val {
                   dp[i+1][1] = val
                   pre[i+1][1] = j
                   sel[i+1][1] = int64(st[i])*B + int64(choose)
               }
           }
       }
   }
   // output result
   result := dp[M+1][0]
   fmt.Fprintln(writer, result)
   // backtrack
   var ans []int64
   j := 0
   for i := M + 1; i >= 1; i-- {
       if sel[i][j] != -1 {
           ans = append(ans, sel[i][j])
       }
       j = pre[i][j]
   }
   K := len(ans)
   fmt.Fprintln(writer, K)
   for _, v := range ans {
       x := v % B
       y := v / B
       // +1 for 1-based
       fmt.Fprintln(writer, x+1, y+1)
   }
}
