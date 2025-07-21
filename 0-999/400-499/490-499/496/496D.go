package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // prefix sums
   p1 := make([]int, n+1)
   p2 := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p1[i] = p1[i-1]
       p2[i] = p2[i-1]
       if a[i] == 1 {
           p1[i]++
       } else {
           p2[i]++
       }
   }
   // positions of each prefix sum value
   max1, max2 := p1[n], p2[n]
   pos1 := make([][]int, max1+1)
   pos2 := make([][]int, max2+1)
   for i := 1; i <= n; i++ {
       pos1[p1[i]] = append(pos1[p1[i]], i)
       pos2[p2[i]] = append(pos2[p2[i]], i)
   }
   type pair struct{ s, t int }
   var res []pair
   INF := n + 5
   // try t
   for t := 1; t <= n; t++ {
       // must have at least one set
       if t > max1 && t > max2 {
           break
       }
       i := 0
       w1, w2 := 0, 0
       lastWin := 0
       valid := true
       for i < n {
           c1, c2 := p1[i], p2[i]
           tgt1, tgt2 := c1+t, c2+t
           next1, next2 := INF, INF
           if tgt1 <= max1 {
               arr := pos1[tgt1]
               j := sort.Search(len(arr), func(j int) bool { return arr[j] > i })
               if j < len(arr) {
                   next1 = arr[j]
               }
           }
           if tgt2 <= max2 {
               arr := pos2[tgt2]
               j := sort.Search(len(arr), func(j int) bool { return arr[j] > i })
               if j < len(arr) {
                   next2 = arr[j]
               }
           }
           j := next1
           winner := 1
           if next2 < next1 {
               j = next2
               winner = 2
           }
           if j == INF {
               valid = false
               break
           }
           if winner == 1 {
               w1++
           } else {
               w2++
           }
           lastWin = winner
           i = j
       }
       if !valid || i != n || w1 == w2 {
           continue
       }
       s := w1
       if w2 > s {
           s = w2
       }
       if (lastWin == 1 && w1 != s) || (lastWin == 2 && w2 != s) {
           continue
       }
       res = append(res, pair{s, t})
   }
   // sort by s, then t
   sort.Slice(res, func(i, j int) bool {
       if res[i].s != res[j].s {
           return res[i].s < res[j].s
       }
       return res[i].t < res[j].t
   })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, len(res))
   for _, pt := range res {
       fmt.Fprintf(writer, "%d %d\n", pt.s, pt.t)
   }
}
