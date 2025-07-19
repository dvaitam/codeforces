package main

import (
   "fmt"
)

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   a := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Scan(&a[i])
   }
   if n == 1 {
       fmt.Println(0)
       return
   }
   // count prefix where a[i] == i
   hm := 1
   for hm <= n && a[hm] == hm {
       hm++
   }
   res := 0
   if hm-2 > res {
       res = hm - 2
   }
   // count suffix where a[j] == 1000 - (n-j)
   hm = n
   co := 1000
   if a[n] == co {
       for hm >= 1 && a[hm] == co {
           hm--
           co--
       }
       tmp := n - hm - 1
       if tmp > res {
           res = tmp
       }
   }
   // check all segments
   for i := 1; i+1 < n; i++ {
       for j := i + 2; j <= n; j++ {
           if a[j]-a[i] == j-i {
               tmp := j - i - 1
               if tmp > res {
                   res = tmp
               }
           }
       }
   }
   fmt.Println(res)
}
