package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MaxN = 5005

var f [][]float64

func init() {
   f = make([][]float64, MaxN)
   for i := 0; i < MaxN; i++ {
       f[i] = make([]float64, i+1)
   }
   for i := 1; i < MaxN; i++ {
       f[i][1] = 1.0
       for j := 2; j <= i; j++ {
           var v1, v2 float64
           if j-2 >= 0 {
               v1 = float64(j-2) * f[i-2][j-2]
           }
           if j-1 <= i-2 {
               v2 = float64(i-j) * f[i-2][j-1]
           }
           f[i][j] = (v1 + v2) / float64(i)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       gOut := make([][]int, n+1)
       gRev := make([][]int, n+1)
       rd := make([]int, n+1)
       for i := 0; i < m; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           gOut[x] = append(gOut[x], y)
           gRev[y] = append(gRev[y], x)
           rd[x]++
       }
       dp := make([]float64, n+1)
       dp[n] = 1.0
       queue := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if rd[i] == 0 {
               queue = append(queue, i)
           }
       }
       for idx := 0; idx < len(queue); idx++ {
           x := queue[idx]
           children := gOut[x]
           now := len(children)
           if now > 0 {
               mp := make([]float64, now)
               for i, child := range children {
                   mp[i] = dp[child]
               }
               sort.Float64s(mp)
               for i := 0; i < now; i++ {
                   dp[x] += f[now][now-i] * mp[i]
               }
           }
           for _, parent := range gRev[x] {
               rd[parent]--
               if rd[parent] == 0 {
                   queue = append(queue, parent)
               }
           }
       }
       fmt.Fprintf(writer, "%.20f\n", dp[1])
   }
}
