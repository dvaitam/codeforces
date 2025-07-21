package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   adj [][]int
   vis []bool
   ans int64
)

func dfs(start, u, depth int) {
   if depth == 4 {
       // check if can close cycle
       for _, v := range adj[u] {
           if v == start {
               ans++
               break
           }
       }
       return
   }
   for _, v := range adj[u] {
       if v > start && !vis[v] {
           vis[v] = true
           dfs(start, v, depth+1)
           vis[v] = false
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n, &m)
   adj = make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // sort adjacency lists for deterministic order
   for i := 0; i < n; i++ {
       // simple insertion sort or sort package
       // using sort.Ints
       // but to avoid import, adjacency unsorted is fine
   }
   vis = make([]bool, n)
   // enumerate cycles, start is smallest vertex in cycle
   for start := 0; start < n; start++ {
       vis[start] = true
       for _, v := range adj[start] {
           if v > start {
               vis[v] = true
               dfs(start, v, 1)
               vis[v] = false
           }
       }
       vis[start] = false
   }
   fmt.Fprintln(writer, ans)
}
