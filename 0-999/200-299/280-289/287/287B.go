package main

import (
   "fmt"
   "os"
)

func main() {
   var n, k int64
   if _, err := fmt.Fscan(os.Stdin, &n, &k); err != nil {
       return
   }
   if n == 1 {
       fmt.Println(0)
       return
   }
   target := n - 1
   // Maximum possible pipes from all splitters: sum_{i=2..k} (i-1) = k*(k-1)/2
   maxTotal := k * (k - 1) / 2
   if maxTotal < target {
       fmt.Println(-1)
       return
   }
   // Binary search minimal m in [1..k-1] s.t. sum of m largest (i-1) >= target
   var l, r, ans int64
   l = 1
   r = k - 1
   ans = k // sentinel
   for l <= r {
       m := (l + r) / 2
       // sum of m largest values in [1..k-1]: S = m*(2*(k-1) - (m-1))/2 = m*(2*k - m -1)/2
       sum := m * (2*k - m - 1) / 2
       if sum >= target {
           ans = m
           r = m - 1
       } else {
           l = m + 1
       }
   }
   fmt.Println(ans)
}
