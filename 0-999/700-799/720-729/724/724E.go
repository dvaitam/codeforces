package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var c int64
   if _, err := fmt.Fscan(reader, &n, &c); err != nil {
       return
   }
   p := make([]int64, n+1)
   s := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p[i])
   }
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // prefix sum of s
   preS := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       preS[i] = preS[i-1] + s[i]
   }
   // suffix sum of p
   sufP := make([]int64, n+2)
   for i := n; i >= 1; i-- {
       sufP[i] = sufP[i+1] + p[i]
   }
   // compute minimum cut cost
   // consider K = number of cities in S prefix, K from 0..n
   // cost = sufP[K+1] + preS[K] + c * K * (n-K)
   var ans int64 = sufP[1] // K=0: sufP[1] + preS[0] + c*0*(n)
   for K := 0; K <= n; K++ {
       // suffix from K+1 to n
       var remP int64 = 0
       if K+1 <= n {
           remP = sufP[K+1]
       }
       sold := preS[K]
       // transport cut cost
       transport := c * int64(K) * int64(n-K)
       cost := remP + sold + transport
       if cost < ans {
           ans = cost
       }
   }
   // ans is the max flow
   fmt.Fprintln(writer, ans)
}
