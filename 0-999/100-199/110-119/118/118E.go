package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents an undirected edge with an identifier
type Edge struct {
   to, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]Edge, n+1)
   for i := 1; i <= m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], Edge{v, i})
       adj[v] = append(adj[v], Edge{u, i})
   }

   visited := make([]bool, n+1)
   dep := make([]int, n+1)
   low := make([]int, n+1)
   parent := make([]int, n+1)
   done := make([]bool, m+1)
   ans := make([][2]int, m+1)
   hasBridge := false

   var dfs func(x int)
   dfs = func(x int) {
       visited[x] = true
       low[x] = dep[x]
       for _, e := range adj[x] {
           j := e.to
           id := e.id
           if !visited[j] {
               if !done[id] {
                   done[id] = true
                   ans[id][0], ans[id][1] = x, j
               }
               parent[j] = x
               dep[j] = dep[x] + 1
               dfs(j)
               if low[j] < low[x] {
                   low[x] = low[j]
               }
               if low[j] > dep[x] {
                   hasBridge = true
               }
           } else if j != parent[x] {
               if dep[j] < low[x] {
                   low[x] = dep[j]
               }
               if !done[id] {
                   done[id] = true
                   ans[id][0], ans[id][1] = x, j
               }
           }
       }
   }

   dfs(1)
   if hasBridge {
       writer.WriteString("0\n")
   } else {
       for i := 1; i <= m; i++ {
           writer.WriteString(fmt.Sprintf("%d %d\n", ans[i][0], ans[i][1]))
       }
   }
}
