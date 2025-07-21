package main

import (
   "bufio"
   "fmt"
   "os"
)

type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := range p {
       p[i] = i
   }
   return &DSU{p}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx != ry {
       d.p[ry] = rx
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   c := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &c[i])
   }
   // map node to type
   typeOf := make([]int, n+1)
   idx := 1
   for i := 1; i <= k; i++ {
       cnt := c[i-1]
       for j := 0; j < cnt; j++ {
           typeOf[idx] = i
           idx++
       }
   }
   // DSU for zero-cost edges
   dsu := NewDSU(n)
   type edge struct{u, v, w int}
   edges := make([]edge, 0, m)
   for i := 0; i < m; i++ {
       var u, v, w int
       fmt.Fscan(reader, &u, &v, &w)
       edges = append(edges, edge{u, v, w})
       if w == 0 {
           dsu.Union(u, v)
       }
   }
   // check each type is zero-connected
   idx = 1
   for i := 1; i <= k; i++ {
       cnt := c[i-1]
       if cnt > 0 {
           root := dsu.Find(idx)
           for j := 0; j < cnt; j++ {
               if dsu.Find(idx+j) != root {
                   fmt.Fprintln(writer, "No")
                   return
               }
           }
       }
       idx += cnt
   }
   // build type graph
   const INF = int(1e9)
   dist := make([][]int, k)
   for i := 0; i < k; i++ {
       dist[i] = make([]int, k)
       for j := 0; j < k; j++ {
           if i == j {
               dist[i][j] = 0
           } else {
               dist[i][j] = INF
           }
       }
   }
   for _, e := range edges {
       if e.w == 0 {
           continue
       }
       ti := typeOf[e.u] - 1
       tj := typeOf[e.v] - 1
       if ti != tj {
           if e.w < dist[ti][tj] {
               dist[ti][tj] = e.w
               dist[tj][ti] = e.w
           }
       }
   }
   // Floyd-Warshall
   for p := 0; p < k; p++ {
       for i := 0; i < k; i++ {
           if dist[i][p] == INF {
               continue
           }
           for j := 0; j < k; j++ {
               nd := dist[i][p] + dist[p][j]
               if nd < dist[i][j] {
                   dist[i][j] = nd
               }
           }
       }
   }
   // output
   fmt.Fprintln(writer, "Yes")
   for i := 0; i < k; i++ {
       for j := 0; j < k; j++ {
           if dist[i][j] >= INF {
               fmt.Fprint(writer, -1)
           } else {
               fmt.Fprint(writer, dist[i][j])
           }
           if j+1 < k {
               fmt.Fprint(writer, " ")
           }
       }
       fmt.Fprintln(writer)
   }
}
