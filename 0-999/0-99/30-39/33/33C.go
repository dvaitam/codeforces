package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // prefix sums P[i] = sum a[1..i]
   P := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       P[i] = P[i-1] + a[i]
   }
   // prefix minimal of P
   minP := make([]int64, n+1)
   const INF = int64(4e18)
   minP[0] = P[0]
   for i := 1; i <= n; i++ {
       if P[i] < minP[i-1] {
           minP[i] = P[i]
       } else {
           minP[i] = minP[i-1]
       }
   }
   // suffix sums Suf[i] = sum a[i..n], with Suf[n+1]=0
   Suf := make([]int64, n+2)
   for i := n; i >= 1; i-- {
       Suf[i] = Suf[i+1] + a[i]
   }
   // total sum
   total := P[n]
   // compute minimal flip cost
   tcost := INF
   // case1: prefix i < suffix j, cost = minP[j-1] + Suf[j]
   for j := 1; j <= n+1; j++ {
       // j==n+1: Suf[n+1]=0, minP[n] covers P[0..n]
       cost := minP[j-1] + Suf[j]
       if cost < tcost {
           tcost = cost
       }
   }
   // case2: prefix i >= suffix j, cost = minP[j-1] + Suf[i+1] for j<=i
   // for fixed i, best j gives minP[i-1]
   for i := 1; i <= n; i++ {
       cost := minP[i-1] + Suf[i+1]
       if cost < tcost {
           tcost = cost
       }
   }
   // maximize sum = total - 2*tcost
   ans := total - 2*tcost
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
