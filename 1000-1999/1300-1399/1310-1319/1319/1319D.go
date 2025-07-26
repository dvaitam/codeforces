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
   fmt.Fscan(reader, &n, &m)
   // forward and reverse graphs
   g := make([][]int, n+1)
   rg := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       g[u] = append(g[u], v)
       rg[v] = append(rg[v], u)
   }
   var k int
   fmt.Fscan(reader, &k)
   p := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &p[i])
   }
   s := p[0]
   t := p[k-1]
   // BFS on reverse graph from t
   const INF = 1e9
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   queue := make([]int, 0, n)
   dist[t] = 0
   queue = append(queue, t)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, w := range rg[u] {
           if dist[w] > dist[u] + 1 {
               dist[w] = dist[u] + 1
               queue = append(queue, w)
           }
       }
   }
   // compute min and max rebuilds
   minR, maxR := 0, 0
   for i := 0; i < k-1; i++ {
       u := p[i]
       v := p[i+1]
       // count of shortest-path outgoing edges from u
       cnt := 0
       for _, w := range g[u] {
           if dist[u] == dist[w] + 1 {
               cnt++
           }
       }
       // check if actual move is shortest
       if dist[u] != dist[v] + 1 {
           minR++
           maxR++
       } else {
           // it is shortest, but if multiple choices, could rebuild
           if cnt > 1 {
               maxR++
           }
       }
   }
   fmt.Fprintln(writer, minR, maxR)
}
