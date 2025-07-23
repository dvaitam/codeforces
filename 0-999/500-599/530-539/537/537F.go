package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   ans := make([]int, n)
   // threshold for brute small k
   const K = 350
   smallK := K
   if smallK > n-1 {
       smallK = n - 1
   }
   // brute for small k
   for k := 1; k <= smallK; k++ {
       cnt := 0
       // for parent j, children at [l, r]
       for j := 1; ; j++ {
           l := k*(j-1) + 2
           if l > n {
               break
           }
           r := k*j + 1
           if r > n {
               r = n
           }
           pj := a[j]
           for v := l; v <= r; v++ {
               if a[v] < pj {
                   cnt++
               }
           }
       }
       ans[k] = cnt
   }
   // difference array for large k
   diff := make([]int, n+2)
   // process u for k > smallK
   for u := 2; u <= n; u++ {
       x := u - 2
       // max parent p for k >= smallK+1: floor(x/(smallK+1))+1
       maxP := x/(smallK+1) + 1
       if maxP > u-1 {
           maxP = u - 1
       }
       for p := 1; p <= maxP; p++ {
           if a[p] <= a[u] {
               continue
           }
           // compute k interval for which parent of u is p
           var L, R int
           if p == 1 {
               L = x + 1
               R = n - 1
           } else {
               L = x/p + 1
               R = x/(p-1)
           }
           // restrict to k > smallK
           if L <= smallK {
               L = smallK + 1
           }
           if L > R || L > n-1 {
               continue
           }
           if R > n-1 {
               R = n - 1
           }
           diff[L]++
           diff[R+1]--
       }
   }
   // accumulate diff into ans for k > smallK
   cur := 0
   for k := smallK + 1; k <= n-1; k++ {
       cur += diff[k]
       ans[k] = cur
   }
   // output
   for k := 1; k <= n-1; k++ {
       if k > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans[k])
   }
   out.WriteByte('\n')
}
