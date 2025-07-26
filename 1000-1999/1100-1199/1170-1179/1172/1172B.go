package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 998244353

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   deg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       deg[u]++
       deg[v]++
   }
   // precompute factorials up to n
   fact := make([]int64, n+1)
   fact[0] = 1
   for i := 1; i <= n; i++ {
       fact[i] = fact[i-1] * int64(i) % mod
   }
   var ans int64 = int64(n) % mod
   for i := 1; i <= n; i++ {
       ans = ans * fact[deg[i]] % mod
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   fmt.Fprint(w, ans)
}
