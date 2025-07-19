package main

import (
   "fmt"
)

func main() {
   var n, m int
   if _, err := fmt.Scan(&n, &m); err != nil {
       return
   }
   // adjacency matrix
   g := make([][]bool, n)
   for i := 0; i < n; i++ {
       g[i] = make([]bool, n)
   }
   deg := make([]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Scan(&u, &v)
       u--
       v--
       if u >= 0 && u < n && v >= 0 && v < n {
           if !g[u][v] {
               g[u][v] = true
               g[v][u] = true
               deg[u]++
               deg[v]++
           }
       }
   }
   // dfs to check connectivity
   mark := make([]bool, n)
   var dfs func(int)
   dfs = func(pos int) {
       mark[pos] = true
       for i := 0; i < n; i++ {
           if !mark[i] && g[pos][i] {
               dfs(i)
           }
       }
   }
   if n > 0 {
       dfs(0)
   }
   // check connected and edge count equals nodes
   for i := 0; i < n; i++ {
       if !mark[i] || m != n {
           fmt.Println("NO")
           return
       }
   }
   // prune leaves
   for rep := 0; rep < n; rep++ {
       for i := 0; i < n; i++ {
           if deg[i] == 1 {
               // find the only neighbor
               var f int
               for f = 0; f < n; f++ {
                   if g[i][f] {
                       break
                   }
               }
               // remove edge
               g[i][f] = false
               g[f][i] = false
               deg[i]--
               deg[f]--
               m--
           }
       }
   }
   // count remaining vertices with degree != 0
   tot := 0
   for i := 0; i < n; i++ {
       if deg[i] != 0 {
           tot++
       }
   }
   if tot == m {
       fmt.Println("FHTAGN!")
   } else {
       fmt.Println("NO")
   }
}
