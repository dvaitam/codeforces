package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   pref := make([]int64, n+1)
   for i := 0; i < n; i++ {
       pref[i+1] = pref[i] + a[i]
   }
   const INF = int64(9e18)
   best := INF
   // only increase scenario
   l := 0
   for r := 0; r < n; r++ {
       for l < r {
           sumInc := int64(r-l)*a[r] - (pref[r] - pref[l])
           if sumInc <= k {
               break
           }
           l++
       }
       length := r - l + 1
       if length >= m {
           start := r - m + 1
           if start < l {
               start = l
           }
           sumSeg := pref[r+1] - pref[start]
           incTime := int64(m)*a[r] - sumSeg
           if incTime < best {
               best = incTime
           }
       }
   }
   // only decrease scenario
   if m <= n {
       Tidx := n - m
       T := a[Tidx]
       sumSeg := pref[n] - pref[Tidx]
       decTime := sumSeg - int64(m)*T
       if decTime < best {
           best = decTime
       }
   }
   if best == INF {
       best = 0
   }
   fmt.Println(best)
}
