package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   a := make([]int64, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Sentinels
   const INF = int64(2e18)
   a[0] = -INF
   a[n+1] = INF
   bpos := make([]int, 0, k+2)
   bpos = append(bpos, 0)
   for i := 0; i < k; i++ {
       var x int
       fmt.Fscan(in, &x)
       bpos = append(bpos, x)
   }
   bpos = append(bpos, n+1)
   res := int64(0)
   // Process segments
   for i := 1; i < len(bpos); i++ {
       l := bpos[i-1]
       r := bpos[i]
       if a[l] >= a[r] {
           fmt.Println(-1)
           return
       }
       // LIS on a[l+1..r-1] with bounds (a[l], a[r])
       d := make([]int64, 0, r-l-1)
       for j := l + 1; j < r; j++ {
           v := a[j]
           if v <= a[l] || v >= a[r] {
               continue
           }
           idx := sort.Search(len(d), func(i int) bool { return d[i] >= v })
           if idx == len(d) {
               d = append(d, v)
           } else {
               d[idx] = v
           }
       }
       segLen := int64(r - l - 1)
       res += segLen - int64(len(d))
   }
   fmt.Println(res)
}
