package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   to, id int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   adj := make([][]edge, n+1)
   for i := 1; i <= m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], edge{to: v, id: i})
       adj[v] = append(adj[v], edge{to: u, id: i})
   }
   visited := make([]bool, n+1)
   depth := make([]int, n+1)
   parent := make([]int, n+1)
   parentEdge := make([]int, n+1)
   good := make([]bool, m+1)
   haveCycle := false

   var dfs func(u int)
   dfs = func(u int) {
       visited[u] = true
       for _, e := range adj[u] {
           v, eid := e.to, e.id
           if eid == parentEdge[u] {
               continue
           }
           if !visited[v] {
               parent[v] = u
               parentEdge[v] = eid
               depth[v] = depth[u] + 1
               dfs(v)
           } else if depth[v] < depth[u] {
               cycleLen := depth[u] - depth[v] + 1
               if cycleLen%2 == 1 {
                   cycleMark := make([]bool, m+1)
                   cycleMark[eid] = true
                   x := u
                   for x != v {
                       cycleMark[parentEdge[x]] = true
                       x = parent[x]
                   }
                   if !haveCycle {
                       for i := 1; i <= m; i++ {
                           good[i] = cycleMark[i]
                       }
                       haveCycle = true
                   } else {
                       for i := 1; i <= m; i++ {
                           good[i] = good[i] && cycleMark[i]
                       }
                   }
               }
           }
       }
   }

   for i := 1; i <= n; i++ {
       if !visited[i] {
           depth[i] = 0
           parent[i] = 0
           parentEdge[i] = 0
           dfs(i)
       }
   }
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   if !haveCycle {
       fmt.Fprintln(out, m)
       for i := 1; i <= m; i++ {
           if i > 1 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, i)
       }
       out.WriteByte('\n')
       return
   }
   var res []int
   for i := 1; i <= m; i++ {
       if good[i] {
           res = append(res, i)
       }
   }
   fmt.Fprintln(out, len(res))
   for i, eid := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, eid)
   }
   out.WriteByte('\n')
}
