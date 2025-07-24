package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
       a[i]-- // zero-index
   }
   vis := make([]bool, n)
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       idx[i] = -1
   }
   var cycles []int
   var stack []int
   for i := 0; i < n; i++ {
       if vis[i] {
           continue
       }
       stack = stack[:0]
       u := i
       for !vis[u] {
           vis[u] = true
           idx[u] = len(stack)
           stack = append(stack, u)
           u = a[u]
       }
       if idx[u] >= 0 {
           // found cycle
           k := len(stack) - idx[u]
           cycles = append(cycles, k)
       }
       // reset idx for nodes in stack
       for _, v := range stack {
           idx[v] = -1
       }
   }
   // precompute powers of 2
   pow2 := make([]int, n+1)
   pow2[0] = 1
   for i := 1; i <= n; i++ {
       pow2[i] = pow2[i-1] * 2 % mod
   }
   totalCycleNodes := 0
   for _, k := range cycles {
       totalCycleNodes += k
   }
   // start with choices for tree edges: 2^(n - totalCycleNodes)
   ans := pow2[n-totalCycleNodes]
   // for each cycle, multiply by (2^k - 2)
   for _, k := range cycles {
       term := pow2[k] - 2
       if term < 0 {
           term += mod
       }
       ans = int(int64(ans) * int64(term) % mod)
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprintln(w, ans)
}
