package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // read tree edges
   type edge struct{ u, v int }
   edges := make([]edge, n-1)
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       edges[i] = edge{u, v}
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // parent and depth via BFS
   parent := make([]int, n)
   depth := make([]int, n)
   parent[0] = -1
   depth[0] = 0
   queue := make([]int, 0, n)
   queue = append(queue, 0)
   for idx := 0; idx < len(queue); idx++ {
       u := queue[idx]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           queue = append(queue, v)
       }
   }
   // values per node (edge to parent)
   value := make([]int, n)
   for i := range value {
       value[i] = 1
   }
   // constraints
   var m int
   fmt.Fscan(in, &m)
   type cons struct{ u, v, w int }
   consList := make([]cons, m)
   for i := 0; i < m; i++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       u--
       v--
       consList[i] = cons{u, v, w}
       // apply update
       uu, vv := u, v
       if depth[uu] < depth[vv] {
           uu, vv = vv, uu
       }
       // lift uu
       for depth[uu] > depth[vv] {
           value[uu] = max(value[uu], w)
           uu = parent[uu]
       }
       // lift both
       for uu != vv {
           value[uu] = max(value[uu], w)
           value[vv] = max(value[vv], w)
           uu = parent[uu]
           vv = parent[vv]
       }
   }
   // verify constraints
   ok := true
   for _, c := range consList {
       u, v, w := c.u, c.v, c.w
       uu, vv := u, v
       if depth[uu] < depth[vv] {
           uu, vv = vv, uu
       }
       found := false
       for depth[uu] > depth[vv] {
           if value[uu] == w {
               found = true
               break
           }
           uu = parent[uu]
       }
       for !found && uu != vv {
           if value[uu] == w || value[vv] == w {
               found = true
               break
           }
           uu = parent[uu]
           vv = parent[vv]
       }
       if !found {
           ok = false
           break
       }
   }
   if !ok {
       fmt.Fprintln(out, -1)
       return
   }
   // prepare answer for edges
   ans := make([]int, n-1)
   for i, e := range edges {
       // child is deeper node
       if parent[e.u] == e.v {
           ans[i] = value[e.u]
       } else {
           ans[i] = value[e.v]
       }
   }
   // output
   for i, w := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, w)
   }
   out.WriteByte('\n')
}
