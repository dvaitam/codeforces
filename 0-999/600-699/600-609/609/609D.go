package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m, k int
   var s int64
   fmt.Fscan(reader, &n, &m, &k, &s)
   orig1 := make([]int64, n+1)
   orig2 := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &orig1[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &orig2[i])
   }
   c1 := make([]int64, n+1)
   c2 := make([]int64, n+1)
   c1day := make([]int, n+1)
   c2day := make([]int, n+1)
   c1[1] = orig1[1]
   c1day[1] = 1
   for i := 2; i <= n; i++ {
       if c1[i-1] > orig1[i] {
           c1[i] = orig1[i]
           c1day[i] = i
       } else {
           c1[i] = c1[i-1]
           c1day[i] = c1day[i-1]
       }
   }
   c2[1] = orig2[1]
   c2day[1] = 1
   for i := 2; i <= n; i++ {
       if c2[i-1] > orig2[i] {
           c2[i] = orig2[i]
           c2day[i] = i
       } else {
           c2[i] = c2[i-1]
           c2day[i] = c2day[i-1]
       }
   }
   type Node struct {
       x int
       c int64
   }
   q1 := make([]Node, 0, m)
   q2 := make([]Node, 0, m)
   for i := 1; i <= m; i++ {
       var t int
       var c int64
       fmt.Fscan(reader, &t, &c)
       if t == 1 {
           q1 = append(q1, Node{i, c})
       } else {
           q2 = append(q2, Node{i, c})
       }
   }
   tot1 := len(q1)
   tot2 := len(q2)
   sort.Slice(q1, func(i, j int) bool { return q1[i].c < q1[j].c })
   sort.Slice(q2, func(i, j int) bool { return q2[i].c < q2[j].c })
   sum1 := make([]int64, tot1+1)
   sum2 := make([]int64, tot2+1)
   for i := 1; i <= tot1; i++ {
       sum1[i] = sum1[i-1] + q1[i-1].c
   }
   for i := 1; i <= tot2; i++ {
       sum2[i] = sum2[i-1] + q2[i-1].c
   }
   jub := func(day int) bool {
       for i := 0; i <= k; i++ {
           if i <= tot1 && k-i <= tot2 {
               cost := sum1[i]*c1[day] + sum2[k-i]*c2[day]
               if cost <= s {
                   return true
               }
           }
       }
       return false
   }
   if !jub(n) {
       fmt.Fprintln(writer, -1)
       return
   }
   l, r := 1, n
   for l < r {
       mid := (l + r) >> 1
       if jub(mid) {
           r = mid
       } else {
           l = mid + 1
       }
   }
   fmt.Fprintln(writer, l)
   for i := 0; i <= k; i++ {
       if i <= tot1 && k-i <= tot2 {
           cost := sum1[i]*c1[l] + sum2[k-i]*c2[l]
           if cost <= s {
               for o := 0; o < i; o++ {
                   fmt.Fprintln(writer, q1[o].x, c1day[l])
               }
               for o := 0; o < k-i; o++ {
                   fmt.Fprintln(writer, q2[o].x, c2day[l])
               }
               return
           }
       }
   }
}
