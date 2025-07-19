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
   fmt.Fscan(reader, &T)
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
       parent := make([]int, n+1)
       depth := make([]int, n+1)
       visited := make([]bool, n+1)
       depth[1] = 1
       visited[1] = true
       parent[1] = 0
       stackU := make([]int, 0, n)
       stackIt := make([]int, 0, n)
       stackU = append(stackU, 1)
       stackIt = append(stackIt, 0)
       p := 1
       for len(stackU) > 0 {
           u := stackU[len(stackU)-1]
           it := stackIt[len(stackIt)-1]
           if it < len(adj[u]) {
               v := adj[u][it]
               stackIt[len(stackIt)-1]++
               if !visited[v] {
                   visited[v] = true
                   parent[v] = u
                   depth[v] = depth[u] + 1
                   if depth[v] > depth[p] {
                       p = v
                   }
                   stackU = append(stackU, v)
                   stackIt = append(stackIt, 0)
               }
           } else {
               stackU = stackU[:len(stackU)-1]
               stackIt = stackIt[:len(stackIt)-1]
           }
       }
       if depth[p] >= (n+1)/2 {
           fmt.Fprintln(writer, "PATH")
           fmt.Fprintln(writer, depth[p])
           u := p
           for u != 0 {
               fmt.Fprint(writer, u, " ")
               u = parent[u]
           }
           fmt.Fprintln(writer)
       } else {
           fmt.Fprintln(writer, "PAIRING")
           maxd := depth[p]
           levels := make([][]int, maxd+1)
           for i := 1; i <= n; i++ {
               if visited[i] {
                   d := depth[i]
                   if d <= maxd {
                       levels[d] = append(levels[d], i)
                   }
               }
           }
           totalPairs := 0
           for d := 1; d <= maxd; d++ {
               totalPairs += len(levels[d]) / 2
           }
           fmt.Fprintln(writer, totalPairs)
           for d := 1; d <= maxd; d++ {
               lst := levels[d]
               for i := 0; i+1 < len(lst); i += 2 {
                   fmt.Fprintln(writer, lst[i], lst[i+1])
               }
           }
       }
   }
}
