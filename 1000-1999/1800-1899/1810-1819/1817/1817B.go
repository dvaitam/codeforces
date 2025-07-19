package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for tt := 0; tt < t; tt++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       adj := make([][]int, n+1)
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       vis := make([]bool, n+1)
       p := make([]int, n+1)
       var target int
       var found bool

       var dfs func(parent, v int)
       dfs = func(parent, v int) {
           vis[v] = true
           p[v] = parent
           for _, nei := range adj[v] {
               if found {
                   return
               }
               if nei == target && target != parent {
                   p[target] = v
                   found = true
                   return
               }
               if !vis[nei] {
                   dfs(v, nei)
                   if found {
                       return
                   }
               }
           }
       }

       var ansEdges [][2]int
       // find vertex with degree >= 4 and cycle
       for i := 1; i <= n; i++ {
           if len(adj[i]) >= 4 {
               // reset
               for j := 1; j <= n; j++ {
                   vis[j] = false
                   p[j] = 0
               }
               target = i
               found = false
               dfs(0, i)
               if found {
                   break
               }
           }
       }
       if !found {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // build fish graph edges
       fmt.Fprintln(writer, "YES")
       inCycle := make([]bool, n+1)
       curr := target
       // record cycle edges
       for p[curr] != target {
           inCycle[curr] = true
           inCycle[p[curr]] = true
           ansEdges = append(ansEdges, [2]int{curr, p[curr]})
           curr = p[curr]
       }
       // last edge to close cycle
       ansEdges = append(ansEdges, [2]int{curr, p[curr]})
       // add two extra edges from target
       cnt := 0
       for _, nei := range adj[target] {
           if cnt >= 2 {
               break
           }
           if !inCycle[nei] {
               ansEdges = append(ansEdges, [2]int{target, nei})
               cnt++
           }
       }
       fmt.Fprintln(writer, len(ansEdges))
       for _, e := range ansEdges {
           fmt.Fprintf(writer, "%d %d\n", e[0], e[1])
       }
   }
}
