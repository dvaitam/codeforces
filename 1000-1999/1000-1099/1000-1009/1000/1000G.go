package main

import (
   "bufio"
   "fmt"
   "os"
)

const LOG = 20

func maxInt64(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   type Edge struct{ to, w int }
   adj := make([][]Edge, n)
   parent := make([]int, n)
   pw := make([]int64, n)
   depth := make([]int, n)
   // read edges and build undirected adj
   for i := 0; i < n-1; i++ {
       var u, v, w int
       fmt.Fscan(reader, &u, &v, &w)
       u--
       v--
       adj[u] = append(adj[u], Edge{v, w})
       adj[v] = append(adj[v], Edge{u, w})
   }
   // BFS to set parent, pw, depth, order
   order := make([]int, 0, n)
   queue := make([]int, 0, n)
   queue = append(queue, 0)
   parent[0] = 0
   pw[0] = 0
   depth[0] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       order = append(order, u)
       for _, e := range adj[u] {
           if u == parent[u] && e.to == u {
               continue
           }
           if e.to == parent[u] {
               continue
           }
           parent[e.to] = u
           pw[e.to] = int64(e.w)
           depth[e.to] = depth[u] + 1
           queue = append(queue, e.to)
       }
   }
   // binary lifting table
   par := make([][]int, LOG)
   par[0] = make([]int, n)
   copy(par[0], parent)
   for i := 1; i < LOG; i++ {
       par[i] = make([]int, n)
       for v := 0; v < n; v++ {
           par[i][v] = par[i-1][par[i-1][v]]
       }
   }
   // dp and dpu
   dp := make([]int64, n)
   dpu := make([]int64, n)
   // bottom-up dp
   for idx := n - 1; idx >= 1; idx-- {
       v := order[idx]
       p := parent[v]
       w := pw[v]
       val := dp[v] + a[v] - 2*w
       if val < 0 {
           val = 0
       }
       dp[p] += val
   }
   // top-down dpu
   dpu[0] = 0
   for _, u := range order {
       if u == 0 {
           continue
       }
       p := parent[u]
       w := pw[u]
       down := dp[u] + a[u] - 2*w
       if down < 0 {
           down = 0
       }
       val := dpu[p] + dp[p] + a[p] - 2*w - down
       if val < 0 {
           val = 0
       }
       dpu[u] = val
   }
   // parsum
   parsum := make([][]int64, LOG)
   parsum[0] = make([]int64, n)
   for _, v := range order {
       if v == 0 {
           parsum[0][v] = 0
           continue
       }
       p := parent[v]
       w := pw[v]
       down := dp[v] + a[v] - 2*w
       if down < 0 {
           down = 0
       }
       parsum[0][v] = dp[p] - down + a[v] - w
   }
   for i := 1; i < LOG; i++ {
       parsum[i] = make([]int64, n)
       for v := 0; v < n; v++ {
           parsum[i][v] = parsum[i-1][v] + parsum[i-1][par[i-1][v]]
       }
   }
   // lca function
   lca := func(x, y int) int {
       if depth[x] < depth[y] {
           x, y = y, x
       }
       diff := depth[x] - depth[y]
       for i := 0; i < LOG; i++ {
           if diff>>i&1 == 1 {
               x = par[i][x]
           }
       }
       if x == y {
           return x
       }
       for i := LOG - 1; i >= 0; i-- {
           if par[i][x] != par[i][y] {
               x = par[i][x]
               y = par[i][y]
           }
       }
       return parent[x]
   }
   // answer queries
   for i := 0; i < q; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       x--
       y--
       g := lca(x, y)
       var ans int64
       for _, u := range []int{x, y} {
           cur := u
           ans += dp[cur]
           diff := depth[cur] - depth[g]
           for i := 0; i < LOG; i++ {
               if diff>>i&1 == 1 {
                   ans += parsum[i][cur]
                   cur = par[i][cur]
               }
           }
       }
       ans -= dp[g]
       ans += dpu[g]
       ans += a[g]
       fmt.Fprintln(writer, ans)
   }
}
