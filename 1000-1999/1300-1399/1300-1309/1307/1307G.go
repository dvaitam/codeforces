package main

import (
   "bufio"
   "fmt"
   "os"
)

// Edge represents a flow edge
type Edge struct {
   to, rev, cap, cost int
}

func addEdge(g [][]Edge, u, v, cap, cost int) {
   g[u] = append(g[u], Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
   g[v] = append(g[v], Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // build graph
   g := make([][]Edge, n)
   for i := 0; i < m; i++ {
       var u, v, c int
       fmt.Fscan(reader, &u, &v, &c)
       // convert to 0-index
       addEdge(g, u-1, v-1, 1, c)
   }
   var q int
   fmt.Fscan(reader, &q)
   x := make([]int64, q)
   ans := make([]float64, q)
   for i := 0; i < q; i++ {
       var xi int64
       fmt.Fscan(reader, &xi)
       x[i] = xi
       ans[i] = 1e18
   }

   const INF = 1 << 60
   sumF := 0
   var sumC int64
   // auxiliary for spfa
   dist := make([]int64, n)
   prevv := make([]int, n)
   preve := make([]int, n)
   inqueue := make([]bool, n)
   qv := make([]int, n+5)

   for {
       // spfa
       for i := 0; i < n; i++ {
           dist[i] = INF
           inqueue[i] = false
       }
       dist[0] = 0
       head, tail := 0, 0
       qv[tail] = 0
       tail++
       inqueue[0] = true
       for head < tail {
           v := qv[head]
           head++
           inqueue[v] = false
           for ei, e := range g[v] {
               if e.cap > 0 && dist[e.to] > dist[v]+int64(e.cost) {
                   dist[e.to] = dist[v] + int64(e.cost)
                   prevv[e.to] = v
                   preve[e.to] = ei
                   if !inqueue[e.to] {
                       qv[tail] = e.to
                       tail++
                       inqueue[e.to] = true
                   }
               }
           }
       }
       if dist[n-1] == INF {
           break
       }
       // flow = 1
       sumC += dist[n-1]
       sumF++
       // update capacities
       v := n - 1
       for v != 0 {
           u := prevv[v]
           ei := preve[v]
           // decrease cap
           g[u][ei].cap--
           // reverse edge
           rev := g[u][ei].rev
           g[v][rev].cap++
           v = u
       }
       // update answers
       for i := 0; i < q; i++ {
           cur := float64(sumC + x[i]) / float64(sumF)
           if cur < ans[i] {
               ans[i] = cur
           }
       }
   }
   // output
   for i := 0; i < q; i++ {
       fmt.Fprintf(writer, "%.12f\n", ans[i])
   }
}
