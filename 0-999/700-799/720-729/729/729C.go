package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   var s, t int64
   fmt.Fscan(reader, &n, &k, &s, &t)
   cars := make([][2]int64, n)
   var maxV int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &cars[i][0], &cars[i][1]) // c, v
       if cars[i][1] > maxV {
           maxV = cars[i][1]
       }
   }
   stations := make([]int64, k+2)
   stations[0] = 0
   for i := 1; i <= k; i++ {
       fmt.Fscan(reader, &stations[i])
   }
   stations[k+1] = s
   sort.Slice(stations, func(i, j int) bool { return stations[i] < stations[j] })
   // compute segments
   m := k + 1
   d := make([]int64, m)
   var maxD int64
   for i := 0; i < m; i++ {
       dist := stations[i+1] - stations[i]
       d[i] = dist
       if dist > maxD {
           maxD = dist
       }
   }
   // time calculation function
   can := func(v int64) bool {
       if v < maxD {
           return false
       }
       var total int64
       for _, dist := range d {
           if v >= 2*dist {
               total += dist
           } else {
               // v >= dist ensured
               total += 3*dist - v
           }
           if total > t {
               return false
           }
       }
       return total <= t
   }
   // binary search minimal v
   l, r := int64(0), maxV+1
   for l < r {
       mid := (l + r) / 2
       if can(mid) {
           r = mid
       } else {
           l = mid + 1
       }
   }
   if l > maxV || !can(l) {
       fmt.Println(-1)
       return
   }
   vReq := l
   // choose minimal cost
   ans := int64(1<<62 - 1)
   for _, cv := range cars {
       if cv[1] >= vReq && cv[0] < ans {
           ans = cv[0]
       }
   }
   if ans == int64(1<<62-1) {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
