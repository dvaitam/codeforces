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
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // initialize graph
       adj := make([][]int, n+1)
       parent := make([]int, n+1)
       for i := 1; i <= n; i++ {
           parent[i] = i
       }
       // union-find functions
       var find func(int) int
       find = func(x int) int {
           if parent[x] != x {
               parent[x] = find(parent[x])
           }
           return parent[x]
       }
       union := func(x, y int) bool {
           fx := find(x)
           fy := find(y)
           if fx != fy {
               parent[fy] = fx
               return true
           }
           return false
       }
       cnt := n
       for i := 0; i < m; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
           if union(u, v) {
               cnt--
           }
       }
       if cnt > 1 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // vis: -1 unknown, 0 excluded, 1 included
       vis := make([]int, n+1)
       for i := 1; i <= n; i++ {
           vis[i] = -1
       }
       vis[1] = 1
       // queue for BFS
       queue := make([]int, 0, n)
       // exclude neighbors of 1
       for _, v := range adj[1] {
           if vis[v] == -1 {
               vis[v] = 0
               queue = append(queue, v)
           }
       }
       // BFS
       for head := 0; head < len(queue); head++ {
           u := queue[head]
           for _, v := range adj[u] {
               if vis[v] == -1 {
                   vis[v] = 1
                   // exclude neighbors of v
                   for _, w := range adj[v] {
                       if vis[w] == -1 {
                           vis[w] = 0
                           queue = append(queue, w)
                       }
                   }
               }
           }
       }
       // collect result
       res := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if vis[i] == 1 {
               res = append(res, i)
           }
       }
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, len(res))
       for i, v := range res {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
