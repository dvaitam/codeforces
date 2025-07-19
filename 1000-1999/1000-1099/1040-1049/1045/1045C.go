package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader *bufio.Reader
   writer *bufio.Writer
   n, m, q, idx, tot int
   e1 [][]int
   e2 [][]int
   dfn, low []int
   sta []int
   top int
   val, cnt, fa, deep, Top, son, size []int
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func Tarjan(u int) {
   idx++
   dfn[u] = idx
   low[u] = idx
   top++
   sta[top] = u
   for _, v := range e1[u] {
       if dfn[v] == 0 {
           Tarjan(v)
           low[u] = min(low[u], low[v])
           if low[v] >= dfn[u] {
               tot++
               for {
                   x := sta[top]
                   top--
                   val[tot] = 1
                   e2[tot] = append(e2[tot], x)
                   e2[x] = append(e2[x], tot)
                   if x == v {
                       break
                   }
               }
               e2[tot] = append(e2[tot], u)
               e2[u] = append(e2[u], tot)
           }
       } else {
           low[u] = min(low[u], dfn[v])
       }
   }
}

func Dfs1(u, f int) {
   size[u] = 1
   for _, v := range e2[u] {
       if v == f {
           continue
       }
       deep[v] = deep[u] + 1
       cnt[v] = cnt[u] + val[v]
       fa[v] = u
       Dfs1(v, u)
       size[u] += size[v]
       if size[v] > size[son[u]] {
           son[u] = v
       }
   }
}

func Dfs2(u, tp, f int) {
   Top[u] = tp
   if son[u] != 0 {
       Dfs2(son[u], tp, u)
   }
   for _, v := range e2[u] {
       if v == f || v == son[u] {
           continue
       }
       Dfs2(v, v, u)
   }
}

func LCA(u, v int) int {
   for Top[u] != Top[v] {
       if deep[Top[u]] > deep[Top[v]] {
           u, v = v, u
       }
       v = fa[Top[v]]
   }
   if deep[u] > deep[v] {
       return v
   }
   return u
}

func Query(u, v int) int {
   ans := deep[u] + deep[v] - cnt[u] - cnt[v]
   lca := LCA(u, v)
   return ans - 2*deep[lca] + cnt[lca] + cnt[fa[lca]]
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m, &q)
   e1 = make([][]int, n+1)
   maxNodes := n + m + 5
   e2 = make([][]int, maxNodes)
   dfn = make([]int, n+1)
   low = make([]int, n+1)
   sta = make([]int, n+1)
   val = make([]int, maxNodes)
   cnt = make([]int, maxNodes)
   fa = make([]int, maxNodes)
   deep = make([]int, maxNodes)
   Top = make([]int, maxNodes)
   son = make([]int, maxNodes)
   size = make([]int, maxNodes)

   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       e1[u] = append(e1[u], v)
       e1[v] = append(e1[v], u)
   }

   idx = 0
   tot = n
   for i := 1; i <= n; i++ {
       if dfn[i] == 0 {
           Tarjan(i)
           cnt[i] = val[i]
           Dfs1(i, 0)
           Dfs2(i, i, 0)
       }
   }

   for i := 0; i < q; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       res := Query(u, v)
       fmt.Fprintln(writer, res)
   }
}
