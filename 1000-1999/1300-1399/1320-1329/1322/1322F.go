package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { l, r int }

var (
   n, m    int
   repr    []int
   dis     []int
   val     []int
   G       [][]int
   fa      [][]int
   ord     []int
   dep     []int
   dp      []int
   typeX   []int
   s       []pair
   reader  = bufio.NewReader(os.Stdin)
   writer  = bufio.NewWriter(os.Stdout)
)

func max(a, b int) int { if a > b { return a } ; return b }
func min(a, b int) int { if a < b { return a } ; return b }

func upd(x *pair, y pair) {
   x.l = max(x.l, y.l)
   x.r = min(x.r, y.r)
}

func find(x int) int {
   if repr[x] != x {
       px := repr[x]
       r := find(px)
       dis[x] ^= dis[px]
       repr[x] = r
   }
   return repr[x]
}

func merge(u, v, d int) {
   fu := find(u)
   fv := find(v)
   if fu != fv {
       repr[fu] = fv
       dis[fu] = d ^ dis[u] ^ dis[v]
   } else if (dis[u]^dis[v]) != d {
       fmt.Fprintln(writer, -1)
       writer.Flush()
       os.Exit(0)
   }
}

func kfa(u, k int) int {
   for i := 0; k > 0; i++ {
       if k&1 != 0 {
           u = fa[u][i]
       }
       k >>= 1
   }
   return u
}

func lca(u, v int) int {
   if dep[u] < dep[v] {
       u, v = v, u
   }
   u = kfa(u, dep[u]-dep[v])
   if u == v {
       return u
   }
   for i := len(fa[u]) - 1; i >= 0; i-- {
       if fa[u][i] != fa[v][i] {
           u = fa[u][i]
           v = fa[v][i]
       }
   }
   return fa[u][0]
}

func solve(u, k int) bool {
   x := pair{1, k}
   // first pass
   for _, v := range G[u] {
       var t *pair
       if find(v) == find(u) {
           t = &x
       } else {
           t = &s[find(v)]
       }
       if dis[v] == dis[u] {
           typeX[v] = 0
           upd(t, pair{dp[v] + 1, k})
       } else {
           typeX[v] = 1
           upd(t, pair{1, (k + 1 - dp[v]) - 1})
       }
   }
   // second pass
   y := pair{1, k}
   for _, v := range G[u] {
       if find(v) != find(u) {
           comp := find(v)
           l, r := s[comp].l, s[comp].r
           if l > k+1-r {
               typeX[v] ^= 1
               // swap to new interval
               nl := k + 1 - r
               nr := k + 1 - l
               s[comp].l = nl
               s[comp].r = nr
               l, r = nl, nr
           }
           upd(&y, pair{l, r})
       }
   }
   a, b := x.l, x.r
   c, d := y.l, y.r
   if max(a, c) <= min(b, d) {
       dp[u] = max(a, c)
   } else {
       // swap y interval
       c, d = k+1-d, k+1-c
       if max(a, c) <= min(b, d) {
           dp[u] = max(a, c)
           for _, v := range G[u] {
               if find(v) != find(u) {
                   typeX[v] ^= 1
               }
           }
       } else {
           return false
       }
   }
   return true
}

func check(k int) bool {
   for i := 1; i <= n; i++ {
       s[i].l = 1
       s[i].r = k
   }
   for i := n; i >= 1; i-- {
       if !solve(ord[i], k) {
           return false
       }
   }
   return true
}

func main() {
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   repr = make([]int, n+1)
   dis = make([]int, n+1)
   val = make([]int, n+1)
   G = make([][]int, n+1)
   fa = make([][]int, n+1)
   ord = make([]int, n+1)
   dep = make([]int, n+1)
   dp = make([]int, n+1)
   typeX = make([]int, n+1)
   s = make([]pair, n+1)
   for i := 1; i <= n; i++ {
       repr[i] = i
       // prepare binary lifting
       fa[i] = make([]int, 20)
   }
   // read tree
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       G[u] = append(G[u], v)
       G[v] = append(G[v], u)
   }
   // dfs pre-order
   stackU := make([]int, 0, n)
   stackP := make([]int, 0, n)
   stackU = append(stackU, 1)
   stackP = append(stackP, 0)
   cnt := 0
   for len(stackU) > 0 {
       u := stackU[len(stackU)-1]
       p := stackP[len(stackP)-1]
       stackU = stackU[:len(stackU)-1]
       stackP = stackP[:len(stackP)-1]
       dep[u] = dep[p] + 1
       cnt++
       ord[cnt] = u
       fa[u][0] = p
       for j := 1; j < 20; j++ {
           fa[u][j] = fa[fa[u][j-1]][j-1]
       }
       // push children in reverse to maintain order
       for i := len(G[u]) - 1; i >= 0; i-- {
           v := G[u][i]
           if v == p {
               continue
           }
           stackU = append(stackU, v)
           stackP = append(stackP, u)
       }
   }
   // read queries
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       l := lca(u, v)
       if u != l {
           val[u]++
           u = kfa(u, dep[u]-dep[l]-1)
           val[u]--
       }
       if v != l {
           val[v]++
           v = kfa(v, dep[v]-dep[l]-1)
           val[v]--
       }
       if u != l && v != l {
           merge(u, v, 1)
       }
   }
   // propagate values
   for i := n; i >= 1; i-- {
       u := ord[i]
       p := fa[u][0]
       val[p] += val[u]
       if val[u] != 0 {
           merge(u, p, 0)
       }
   }
   // binary search
   l, r, ans := 1, n, 0
   for l <= r {
       mid := (l + r) / 2
       if check(mid) {
           ans = mid
           r = mid - 1
       } else {
           l = mid + 1
       }
   }
   // final assignment
   check(ans)
   fmt.Fprintln(writer, ans)
   for i := 1; i <= n; i++ {
       u := ord[i]
       typeX[u] ^= typeX[fa[u][0]]
       // print per original index order
       // store result in dp when typeX==0 else ans+1-dp
       var res int
       if typeX[u] == 0 {
           res = dp[u]
       } else {
           res = ans + 1 - dp[u]
       }
       fmt.Fprint(writer, res)
       if i < n {
           fmt.Fprint(writer, ' ')
       }
   }
}
