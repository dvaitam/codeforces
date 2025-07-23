package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

var (
   n      int
   adj    [][]int
   dpDown []int64
   dpUp   []int64
   ans    []int64
)

func dfs1(u, p int) {
   var prod int64 = 1
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dfs1(v, u)
       prod = prod * (dpDown[v] + 1) % mod
   }
   dpDown[u] = prod
}

func dfs2(u, p int) {
   deg := len(adj[u])
   f := make([]int64, deg)
   for i, v := range adj[u] {
       if v == p {
           f[i] = dpUp[u] + 1
       } else {
           f[i] = dpDown[v] + 1
       }
   }
   pre := make([]int64, deg+1)
   suf := make([]int64, deg+1)
   pre[0] = 1
   for i := 0; i < deg; i++ {
       pre[i+1] = pre[i] * f[i] % mod
   }
   suf[deg] = 1
   for i := deg - 1; i >= 0; i-- {
       suf[i] = suf[i+1] * f[i] % mod
   }
   ans[u] = pre[deg]
   for i, v := range adj[u] {
       if v == p {
           continue
       }
       // excluding f[i]
       dpUp[v] = pre[i] * suf[i+1] % mod
       dfs2(v, u)
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // read n
   fmt.Fscan(reader, &n)
   adj = make([][]int, n+1)
   // parent pointers p2..pn
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       adj[p] = append(adj[p], i)
       adj[i] = append(adj[i], p)
   }
   dpDown = make([]int64, n+1)
   dpUp = make([]int64, n+1)
   ans = make([]int64, n+1)
   // root at 1
   dpUp[1] = 0
   dfs1(1, 0)
   dfs2(1, 0)
   // output
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", ans[i]))
   }
   writer.WriteByte('\n')
}
