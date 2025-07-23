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
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   visited := make([]bool, n)
   color := make([]int, n)
   total := 0
   queue := make([]int, 0, n)
   dist := make([]int, n)
   for i := 0; i < n; i++ {
       if visited[i] {
           continue
       }
       // BFS for bipartiteness and collect component nodes
       comp := make([]int, 0)
       queue = queue[:0]
       visited[i] = true
       color[i] = 0
       queue = append(queue, i)
       comp = append(comp, i)
       ok := true
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if !visited[v] {
                   visited[v] = true
                   color[v] = 1 - color[u]
                   queue = append(queue, v)
                   comp = append(comp, v)
               } else if color[v] == color[u] {
                   ok = false
                   break
               }
           }
           if !ok {
               break
           }
       }
       if !ok {
           fmt.Fprintln(writer, -1)
           return
       }
       // BFS to find furthest node from comp[0]
       for _, u := range comp {
           dist[u] = -1
       }
       start := comp[0]
       dist[start] = 0
       queue = queue[:0]
       queue = append(queue, start)
       var far int = start
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if dist[v] == -1 {
                   dist[v] = dist[u] + 1
                   queue = append(queue, v)
               }
           }
           if dist[u] > dist[far] {
               far = u
           }
       }
       // BFS from far to get diameter
       for _, u := range comp {
           dist[u] = -1
       }
       dist[far] = 0
       queue = queue[:0]
       queue = append(queue, far)
       diam := 0
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if dist[v] == -1 {
                   dist[v] = dist[u] + 1
                   queue = append(queue, v)
               }
           }
           if dist[u] > diam {
               diam = dist[u]
           }
       }
       total += diam
   }
   fmt.Fprintln(writer, total)
}
