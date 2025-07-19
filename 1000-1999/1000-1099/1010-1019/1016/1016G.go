package main

import (
   "bufio"
   "fmt"
   "os"
)

const LG = 20

type Edge struct { to int; w int64 }

var (
   n, q int
   v    [][]Edge
   a    []int64
   dp, dpu []int64
   par   [][]int
   parsum [][]int64
   tin, tout, h []int
   timer int
   outBuf *bufio.Writer
)

func dfs1(u, p, depth int) {
   h[u] = depth
   timer++
   tin[u] = timer
   par[u][0] = p
   for i := 1; i < LG; i++ {
       par[u][i] = par[par[u][i-1]][i-1]
   }
   for _, e := range v[u] {
       if e.to == p {
           continue
       }
       dfs1(e.to, u, depth+1)
       gain := dp[e.to] + a[e.to] - 2*e.w
       if gain > 0 {
           dp[u] += gain
       }
   }
   tout[u] = timer
}

func dfs2(u, p int, d int64) {
   if d > 0 {
       dpu[u] = d
   }
   for _, e := range v[u] {
       if e.to == p {
           continue
       }
       gain := dp[e.to] + a[e.to] - 2*e.w
       if gain < 0 {
           gain = 0
       }
       // exclude child's contribution and add others
       nd := dpu[u] + dp[u] + a[u] - 2*e.w - gain
       dfs2(e.to, u, nd)
   }
}

func dfs3(u, p int, pw int64) {
   var upGain int64
   if u == p {
       upGain = 0
   } else {
       childGain := dp[u] + a[u] - 2*pw
       if childGain < 0 {
           childGain = 0
       }
       upGain = dp[p] - childGain + a[u] - pw
   }
   parsum[u][0] = upGain
   for i := 1; i < LG; i++ {
       parsum[u][i] = parsum[u][i-1] + parsum[par[u][i-1]][i-1]
   }
   for _, e := range v[u] {
       if e.to == p {
           continue
       }
       dfs3(e.to, u, e.w)
   }
}

func isPar(x, y int) bool {
   return tin[x] <= tin[y] && tout[y] <= tout[x]
}

func lca(x, y int) int {
   if isPar(x, y) {
       return x
   }
   if isPar(y, x) {
       return y
   }
   for i := LG-1; i >= 0; i-- {
       if !isPar(par[x][i], y) {
           x = par[x][i]
       }
   }
   return par[x][0]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   outBuf = bufio.NewWriter(os.Stdout)
   defer outBuf.Flush()
   fmt.Fscan(in, &n, &q)
   a = make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   v = make([][]Edge, n)
   for i := 0; i < n-1; i++ {
       var x, y int
       var w int64
       fmt.Fscan(in, &x, &y, &w)
       x--; y--
       v[x] = append(v[x], Edge{y, w})
       v[y] = append(v[y], Edge{x, w})
   }
   dp = make([]int64, n)
   dpu = make([]int64, n)
   par = make([][]int, n)
   parsum = make([][]int64, n)
   tin = make([]int, n)
   tout = make([]int, n)
   h = make([]int, n)
   for i := 0; i < n; i++ {
       par[i] = make([]int, LG)
       parsum[i] = make([]int64, LG)
   }
   // root at 0
   timer = 0
   dfs1(0, 0, 0)
   dfs2(0, 0, 0)
   dfs3(0, 0, 0)
   for i := 0; i < q; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       x--; y--
       g := lca(x, y)
       var ans int64
       for _, u := range []int{x, y} {
           cur := u
           ans += dp[cur]
           dif := h[cur] - h[g]
           for i := 0; i < LG; i++ {
               if dif&(1<<i) != 0 {
                   ans += parsum[cur][i]
                   cur = par[cur][i]
               }
           }
       }
       ans -= dp[g]
       ans += dpu[g]
       ans += a[g]
       fmt.Fprintln(outBuf, ans)
   }
}
