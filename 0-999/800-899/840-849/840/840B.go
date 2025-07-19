package main

import (
   "bufio"
   "fmt"
   "os"
)

type edge struct {
   v, id int
}

var (
   reader *bufio.Reader
   writer *bufio.Writer
)

func init() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
}

func main() {
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   d := make([]int, n+1)
   x, s := 0, 0
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &d[i])
       if d[i] == -1 {
           x = i
       } else {
           s += d[i]
       }
   }
   if x == 0 && (s&1) == 1 {
       fmt.Fprintln(writer, -1)
       return
   }
   adj := make([][]edge, n+1)
   for i := 1; i <= m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], edge{v: v, id: i})
       adj[v] = append(adj[v], edge{v: u, id: i})
   }
   vis := make([]bool, n+1)
   use := make([]bool, m+1)
   var dfs func(u, e int)
   dfs = func(u, e int) {
       vis[u] = true
       use[e] = (d[u] == 1)
       for _, ed := range adj[u] {
           v, id := ed.v, ed.id
           if !vis[v] {
               dfs(v, id)
               use[e] = use[e] != use[id]
           }
       }
   }
   start := 1
   if x != 0 {
       start = x
   }
   dfs(start, 0)
   cnt := 0
   for i := 1; i <= m; i++ {
       if use[i] {
           cnt++
       }
   }
   fmt.Fprintln(writer, cnt)
   for i := 1; i <= m; i++ {
       if use[i] {
           fmt.Fprintln(writer, i)
       }
   }
}
