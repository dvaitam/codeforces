package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // INF value for unreachable distances, must exceed max possible path length
   INF := n + 5
   // distance matrix
   g := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       g[i] = make([]int, n+1)
       for j := 0; j <= n; j++ {
           if i == j {
               g[i][j] = 0
           } else {
               g[i][j] = INF
           }
       }
   }
   // adjacency
   G := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       G[u] = append(G[u], v)
       G[v] = append(G[v], u)
       g[u][v] = 1
       g[v][u] = 1
   }
   // Floyd-Warshall
   for k := 1; k <= n; k++ {
       for i := 1; i <= n; i++ {
           gik := g[i][k]
           for j := 1; j <= n; j++ {
               if gik+g[k][j] < g[i][j] {
                   g[i][j] = gik + g[k][j]
               }
           }
       }
   }
   // prepare helpers
   vis1 := make([]int, n+1)
   vis2 := make([]int, n+1)
   P := make([]float64, n+1)
   Maxp := make([]float64, n+1)
   tag1, tag2 := 0, 0
   Ans := 0.0
   // buckets by distance
   d := make([][]int, n)
   for i := 0; i < n; i++ {
       d[i] = make([]int, 0, 4)
   }
   po := make([]int, 0, n)
   Le := make([]int, 0, n)
   // main loops
   for i := 1; i <= n; i++ {
       p := 0.0
       // clear buckets
       for l := 0; l < n; l++ {
           d[l] = d[l][:0]
       }
       for j := 1; j <= n; j++ {
           dist := g[i][j]
           if dist < n {
               d[dist] = append(d[dist], j)
           }
       }
       for l := 0; l < n; l++ {
           sz := len(d[l])
           if sz == 0 {
               continue
           }
           if sz == 1 {
               p += 1.0 / float64(n)
               continue
           }
           // collect probability contributions
           tag1++
           po = po[:0]
           invSz := 1.0 / float64(sz)
           for _, v := range d[l] {
               invDeg := invSz / float64(len(G[v]))
               for _, pos := range G[v] {
                   if vis1[pos] != tag1 {
                       vis1[pos] = tag1
                       po = append(po, pos)
                       P[pos] = 0
                   }
                   P[pos] += invDeg
               }
           }
           Max := 0.0
           // for each potential meeting point j
           for j := 1; j <= n; j++ {
               tag2++
               Le = Le[:0]
               q := 0.0
               for _, v := range po {
                   dvg := g[v][j]
                   if vis2[dvg] != tag2 {
                       vis2[dvg] = tag2
                       Le = append(Le, dvg)
                       Maxp[dvg] = 0
                   }
                   if P[v] > Maxp[dvg] {
                       Maxp[dvg] = P[v]
                   }
               }
               for _, dv := range Le {
                   q += Maxp[dv]
               }
               if q > Max {
                   Max = q
               }
           }
           p += mathMax(1.0/float64(n), Max*float64(sz)/float64(n))
       }
       if p > Ans {
           Ans = p
       }
   }
   fmt.Printf("%.10f", Ans)
}

func mathMax(a, b float64) float64 {
   if a > b {
       return a
   }
   return b
}
