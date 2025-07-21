package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func min(a, b int) int { if a < b { return a }; return b }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // adjacency list
   adj := make([][]edge, n+1)
   type EdgeInfo struct{u, v, w int}
   edges := make([]EdgeInfo, 0, n-1)
   for i := 0; i < n-1; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       adj[u] = append(adj[u], edge{to: v, w: w})
       adj[v] = append(adj[v], edge{to: u, w: w})
       edges = append(edges, EdgeInfo{u, v, w})
   }
   // sort adjacency for deterministic DFS order
   for u := 1; u <= n; u++ {
       sort.Slice(adj[u], func(i, j int) bool { return adj[u][i].to < adj[u][j].to })
   }
   // parent and sizes
   parent := make([]int, n+1)
   sizes := make([]int, n+1)
   order := make([]int, 0, n)
   type st struct{ u, idx int }
   stack := make([]st, 0, n)
   parent[1] = 0
   stack = append(stack, st{1, 0})
   // iterative DFS for pre-order and subtree sizes
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       u := top.u
       if top.idx == 0 {
           order = append(order, u)
       }
       if top.idx < len(adj[u]) {
           v := adj[u][top.idx].to
           stack[len(stack)-1].idx++
           if v == parent[u] {
               continue
           }
           parent[v] = u
           stack = append(stack, st{v, 0})
       } else {
           // compute subtree size
           sz := 1
           for _, e := range adj[u] {
               v := e.to
               if v == parent[u] {
                   continue
               }
               sz += sizes[v]
           }
           sizes[u] = sz
           stack = stack[:len(stack)-1]
       }
   }
   // compute maximum sum: sum over edges of 2 * w * min(subtree, n - subtree)
   var total int64
   for _, e := range edges {
       u, v, w := e.u, e.v, e.w
       var s int
       if parent[u] == v {
           s = sizes[u]
       } else if parent[v] == u {
           s = sizes[v]
       } else {
           // should not happen
           continue
       }
       m := min(s, n-s)
       total += int64(w) * 2 * int64(m)
   }
   // build permutation by rotating DFS order by k = floor(n/2)
   k := n / 2
   p := make([]int, n+1)
   for i, u := range order {
       j := (i + k) % n
       p[u] = order[j]
   }
   // output
   fmt.Fprintln(out, total)
   for i := 1; i <= n; i++ {
       if i > 1 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, p[i])
   }
   out.WriteByte('\n')
}

type edge struct{
   to, w int
}
