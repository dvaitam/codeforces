package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s int64
   // read n and s
   fmt.Fscan(reader, &n, &s)
   a := make([]int64, n)
   var total int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       total += a[i]
   }
   if total < s {
       fmt.Println(-1)
       return
   }
   var lo, hi int64 = 0, 1e9
   var ans int64 = 0
   for lo <= hi {
       mid := (lo + hi) >> 1
       var t int64
       mn := int64(1<<63 - 1) // max int64
       for _, v := range a {
           if v >= mid {
               t += v - mid
               if mn > mid {
                   mn = mid
               }
           } else {
               if mn > v {
                   mn = v
               }
           }
       }
       if t >= s {
           // mid valid, try higher
           lo = mid + 1
           if mn > ans {
               ans = mn
           }
       } else {
           hi = mid - 1
       }
   }
   fmt.Println(ans)
}
