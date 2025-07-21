package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var l, r, Ql, Qr int64
   fmt.Fscan(in, &n, &l, &r, &Ql, &Qr)
   w := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &w[i])
   }
   pre := make([]int64, n+1)
   for i := 0; i < n; i++ {
       pre[i+1] = pre[i] + w[i]
   }
   total := pre[n]
   ans := int64(9e18)
   for L := 0; L <= n; L++ {
       leftCost := pre[L] * l
       rightCost := (total - pre[L]) * r
       var penalty int64
       R := n - L
       if L > R {
           diff := int64(L - R - 1)
           if diff > 0 {
               penalty = diff * Ql
           }
       } else if R > L {
           diff := int64(R - L - 1)
           if diff > 0 {
               penalty = diff * Qr
           }
       }
       cost := leftCost + rightCost + penalty
       if cost < ans {
           ans = cost
       }
   }
   fmt.Println(ans)
}
