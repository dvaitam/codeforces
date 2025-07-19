package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m   int
   g       [][]int
   fa, deep, son, top []int
   siz     []int
   sum, src []float64
)

func dfs1(x int) {
   siz[x] = 1
   deep[x] = deep[fa[x]] + 1
   for _, to := range g[x] {
       if to != fa[x] {
           fa[to] = x
           dfs1(to)
           siz[x] += siz[to]
           sum[x] += sum[to] + float64(siz[to])
           if siz[son[x]] < siz[to] {
               son[x] = to
           }
       }
   }
}

func dfs2(x int) {
   if x == son[fa[x]] {
       top[x] = top[fa[x]]
   } else {
       top[x] = x
   }
   for _, to := range g[x] {
       if fa[to] == x {
           dfs2(to)
       }
   }
}

func dfs3(x int) {
   src[x] += sum[x]
   for _, to := range g[x] {
       if fa[to] == x {
           src[to] += src[x] - sum[to] - float64(siz[to]) + float64(n - siz[to])
           dfs3(to)
       }
   }
}

func lca(u, v int) int {
   for top[u] != top[v] {
       if deep[top[u]] > deep[top[v]] {
           u = fa[top[u]]
       } else {
           v = fa[top[v]]
       }
   }
   if deep[u] < deep[v] {
       return u
   }
   return v
}

func solve(u, f int) int {
   for {
       if deep[fa[top[u]]] > deep[f] {
           u = fa[top[u]]
       } else if fa[u] == f {
           return u
       } else if deep[top[u]] > deep[f] {
           u = top[u]
       } else {
           return son[f]
       }
   }
}

func ask(u, v int) float64 {
   l := lca(u, v)
   if u != l && v != l {
       return sum[u]/float64(siz[u]) + sum[v]/float64(siz[v]) + float64(deep[u]+deep[v]+1-2*deep[l])
   }
   if u == l {
       u, v = v, u
   }
   nw := solve(u, v)
   x := sum[u] / float64(siz[u])
   y := (src[v] - sum[nw] - float64(siz[nw])) / float64(n - siz[nw])
   return x + y + float64(deep[u]-deep[v]+1)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m)
   g = make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   fa = make([]int, n+1)
   deep = make([]int, n+1)
   son = make([]int, n+1)
   top = make([]int, n+1)
   siz = make([]int, n+1)
   sum = make([]float64, n+1)
   src = make([]float64, n+1)
   // root at 1
   deep[0] = 0
   fa[1] = 0
   dfs1(1)
   dfs2(1)
   dfs3(1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       ans := ask(u, v)
       fmt.Fprintf(writer, "%.7f\n", ans)
   }
}
