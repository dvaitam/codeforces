package main

import (
   "fmt"
   "os"
)

func main() {
   var n int
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(os.Stdin, &a[i])
   }
   // Remove one hold (except first and last) to minimize max adjacent difference
   const INF = 1000000000
   ans := INF
   // try removing hold at index i (1-based second hold has i=1 in 0-indexed)
   for rem := 1; rem < n-1; rem++ {
       maxd := 0
       prev := a[0]
       for j := 1; j < n; j++ {
           if j == rem {
               continue
           }
           d := a[j] - prev
           if d > maxd {
               maxd = d
           }
           prev = a[j]
       }
       if maxd < ans {
           ans = maxd
       }
   }
   fmt.Fprintln(os.Stdout, ans)
}
