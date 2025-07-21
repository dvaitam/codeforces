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
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS to compute parent[0], depth, order
   const LOG = 18
   parent := make([][]int, LOG)
   for i := range parent {
       parent[i] = make([]int, n+1)
   }
   depth := make([]int, n+1)
   order := make([]int, 0, n)
   // queue
   q := make([]int, 0, n)
   q = append(q, 1)
   parent[0][1] = 0
   depth[1] = 0
   for idx := 0; idx < len(q); idx++ {
       u := q[idx]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[0][u] {
               continue
           }
           parent[0][v] = u
           depth[v] = depth[u] + 1
           q = append(q, v)
       }
   }
   // subtree sizes
   sz := make([]int, n+1)
   for i := 1; i <= n; i++ {
       sz[i] = 1
   }
   for i := len(order) - 1; i > 0; i-- {
       u := order[i]
       p := parent[0][u]
       sz[p] += sz[u]
   }
   // binary lifting
   for k := 1; k < LOG; k++ {
       for v := 1; v <= n; v++ {
           parent[k][v] = parent[k-1][ parent[k-1][v] ]
       }
   }
   // LCA and jump funcs
   var lca func(int,int) int
   lca = func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       // lift u
       dd := depth[u] - depth[v]
       for k := 0; k < LOG; k++ {
           if dd>>k & 1 == 1 {
               u = parent[k][u]
           }
       }
       if u == v {
           return u
       }
       for k := LOG-1; k >= 0; k-- {
           if parent[k][u] != parent[k][v] {
               u = parent[k][u]
               v = parent[k][v]
           }
       }
       return parent[0][u]
   }
   jump := func(u, d int) int {
       for k := 0; k < LOG && u != 0; k++ {
           if d>>k & 1 == 1 {
               u = parent[k][u]
           }
       }
       return u
   }
   var m int
   fmt.Fscan(in, &m)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       if x == y {
           fmt.Fprintln(out, n)
           continue
       }
       w := lca(x, y)
       d := depth[x] + depth[y] - 2*depth[w]
       if d & 1 == 1 {
           fmt.Fprintln(out, 0)
           continue
       }
       // find center c at distance d/2 from x towards y
       k := d / 2
       var c int
       dx := depth[x] - depth[w]
       if k <= dx {
           c = jump(x, k)
       } else {
           // down from w towards y
           rem := k - dx
           // steps from y up to center: (depth[y]-depth[w]) - rem
           c = jump(y, depth[y]-depth[w]-rem)
       }
       // neighbor towards x
       var nx int
       if lca(c, x) == c {
           // c is ancestor of x, take child
           nx = jump(x, depth[x]-depth[c]-1)
       } else {
           nx = parent[0][c]
       }
       // neighbor towards y
       var ny int
       if lca(c, y) == c {
           ny = jump(y, depth[y]-depth[c]-1)
       } else {
           ny = parent[0][c]
       }
       // sizes
       var sx, sy int
       if parent[0][nx] == c {
           sx = sz[nx]
       } else {
           sx = n - sz[c]
       }
       if parent[0][ny] == c {
           sy = sz[ny]
       } else {
           sy = n - sz[c]
       }
       ans := n - sx - sy
       fmt.Fprintln(out, ans)
   }
}
