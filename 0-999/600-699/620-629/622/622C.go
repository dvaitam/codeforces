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

   var n, m int
   fmt.Fscan(in, &n, &m)
   a := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // next different element index
   nxt := make([]int, n+2)
   nxt[n] = n + 1
   // a[n+1] is zero, different from any 1<=a[i]<=1e6
   for i := n - 1; i >= 1; i-- {
       if a[i] == a[i+1] {
           nxt[i] = nxt[i+1]
       } else {
           nxt[i] = i + 1
       }
   }
   // process queries
   for ; m > 0; m-- {
       var l, r, x int
       fmt.Fscan(in, &l, &r, &x)
       if a[l] != x {
           fmt.Fprintln(out, l)
       } else if nxt[l] <= r {
           fmt.Fprintln(out, nxt[l])
       } else {
           fmt.Fprintln(out, -1)
       }
   }
}
