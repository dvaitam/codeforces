package main

import (
   "bufio"
   "fmt"
   "os"
)

type Edge struct {
   to, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   adj := make([][]Edge, n+1)
   // store parent edge id for each node
   edges := make([][2]int, n)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], Edge{v, i})
       adj[v] = append(adj[v], Edge{u, i})
       edges[i] = [2]int{u, v}
   }
   // BFS to get parent, depth, order
   parent := make([]int, n+1)
   parentEdge := make([]int, n+1)
   depth := make([]int, n+1)
   order := make([]int, 0, n)
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   parent[1] = 0
   depth[1] = 0
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       order = append(order, u)
       for _, e := range adj[u] {
           v := e.to
           if v == parent[u] {
               continue
           }
           parent[v] = u
           parentEdge[v] = e.id
           depth[v] = depth[u] + 1
           queue = append(queue, v)
       }
   }
   // binary lifting ancestors
   const maxLog = 18
   p := make([][]int, maxLog)
   p[0] = make([]int, n+1)
   for v := 1; v <= n; v++ {
       p[0][v] = parent[v]
   }
   for k := 1; k < maxLog; k++ {
       p[k] = make([]int, n+1)
       for v := 1; v <= n; v++ {
           p[k][v] = p[k-1][p[k-1][v]]
       }
   }
   // lca function
   lca := func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       diff := depth[u] - depth[v]
       for k := 0; k < maxLog; k++ {
           if diff&(1<<k) != 0 {
               u = p[k][u]
           }
       }
       if u == v {
           return u
       }
       for k := maxLog - 1; k >= 0; k-- {
           if p[k][u] != p[k][v] {
               u = p[k][u]
               v = p[k][v]
           }
       }
       return parent[u]
   }
   // read queries and mark
   var k int
   fmt.Fscan(reader, &k)
   cnt := make([]int64, n+1)
   for i := 0; i < k; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       cnt[a]++
       cnt[b]++
       l := lca(a, b)
       cnt[l] -= 2
   }
   // accumulate counts from leaves to root
   ans := make([]int64, n)
   for i := len(order) - 1; i > 0; i-- {
       v := order[i]
       eid := parentEdge[v]
       ans[eid] = cnt[v]
       cnt[parent[v]] += cnt[v]
   }
   // output answers
   for i := 1; i < n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", ans[i])
   }
   writer.WriteByte('\n')
}
