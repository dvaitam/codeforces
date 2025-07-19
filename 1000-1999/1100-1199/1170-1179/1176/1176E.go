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
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       adj := make([][]int, n+1)
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       dep := make([]int, n+1)
       vis := make([]bool, n+1)
       queue := make([]int, 0, n)
       // BFS from node 1
       queue = append(queue, 1)
       vis[1] = true
       for i := 0; i < len(queue); i++ {
           u := queue[i]
           for _, v := range adj[u] {
               if !vis[v] {
                   vis[v] = true
                   dep[v] = dep[u] + 1
                   queue = append(queue, v)
               }
           }
       }
       groups := [2][]int{}
       for i := 1; i <= n; i++ {
           p := dep[i] & 1
           groups[p] = append(groups[p], i)
       }
       // choose smaller group
       if len(groups[0]) <= len(groups[1]) {
           fmt.Fprintln(writer, len(groups[0]))
           for _, x := range groups[0] {
               fmt.Fprint(writer, x, " ")
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, len(groups[1]))
           for _, x := range groups[1] {
               fmt.Fprint(writer, x, " ")
           }
           fmt.Fprintln(writer)
       }
   }
}
