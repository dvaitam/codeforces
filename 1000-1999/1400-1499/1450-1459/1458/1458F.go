package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXLOG = 17

var (
   n     int
   adj   [][]int
   up    [][]int
   depth []int
)

func dfs(u, p int) {
   up[u][0] = p
   for k := 1; k <= MAXLOG; k++ {
       up[u][k] = up[up[u][k-1]][k-1]
   }
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       depth[v] = depth[u] + 1
       dfs(v, u)
   }
}

func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   diff := depth[u] - depth[v]
   for k := 0; k <= MAXLOG; k++ {
       if diff&(1<<k) != 0 {
           u = up[u][k]
       }
   }
   if u == v {
       return u
   }
   for k := MAXLOG; k >= 0; k-- {
       if up[u][k] != up[v][k] {
           u = up[u][k]
           v = up[v][k]
       }
   }
   return up[u][0]
}

func dist(u, v int) int {
   w := lca(u, v)
   return depth[u] + depth[v] - 2*depth[w]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   up = make([][]int, n+1)
   for i := 0; i <= n; i++ {
       up[i] = make([]int, MAXLOG+1)
   }
   depth = make([]int, n+1)
   depth[1] = 0
   dfs(1, 0)

   var ans int64
   for l := 1; l <= n; l++ {
       a, b := l, l
       diam := 0
       for r := l; r <= n; r++ {
           if r > l {
               d1 := dist(a, r)
               if d1 > diam {
                   diam = d1
                   b = r
               } else {
                   d2 := dist(b, r)
                   if d2 > diam {
                       diam = d2
                       a = r
                   }
               }
           }
           ans += int64(diam)
       }
   }
   fmt.Fprintln(writer, ans)
}
