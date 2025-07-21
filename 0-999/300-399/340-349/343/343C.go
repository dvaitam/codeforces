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
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   h := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &h[i])
   }
   p := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &p[i])
   }
   var lo, hi int64 = 0, 40000000000
   for lo < hi {
       mid := (lo + hi) / 2
       if canCover(h, p, mid) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(lo)
}

func canCover(h, p []int64, T int64) bool {
   n, m := len(h), len(p)
   j := 0
   for i := 0; i < n && j < m; i++ {
       hi := h[i]
       if p[j] > hi {
           limit := hi + T
           // first index > limit
           t := sort.Search(m, func(x int) bool { return p[x] > limit })
           j = t
       } else {
           l := p[j]
           if hi - l > T {
               return false
           }
           k := sort.Search(m, func(x int) bool { return p[x] > hi }) - 1
           // go left then right: r <= T + 2*l - hi
           tMax := k
           t1lim := T + 2*l - hi
           if t1lim >= l {
               idx := sort.Search(m, func(x int) bool { return p[x] > t1lim }) - 1
               if idx > tMax {
                   tMax = idx
               }
           }
           rem := T - (hi - l)
           if rem >= 0 {
               t2lim := hi + rem/2
               idx2 := sort.Search(m, func(x int) bool { return p[x] > t2lim }) - 1
               if idx2 > tMax {
                   tMax = idx2
               }
           }
           j = tMax + 1
       }
   }
   return j >= m
}
