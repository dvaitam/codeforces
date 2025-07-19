package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

// Edge represents a pair of points with squared distance d
type Edge struct {
   d  int
   u, v int
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   x := make([]int, n)
   y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rdr, &x[i], &y[i])
   }
   // prepare edges
   tot := n*(n-1)/2
   edges := make([]Edge, 0, tot)
   for i := 0; i < n; i++ {
       xi, yi := x[i], y[i]
       for j := 0; j < i; j++ {
           dx := xi - x[j]
           dy := yi - y[j]
           d2 := dx*dx + dy*dy
           edges = append(edges, Edge{d: d2, u: i, v: j})
       }
   }
   // sort descending by distance
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].d > edges[j].d
   })
   // bitsets per vertex
   g := (n + 63) / 64
   bits := make([][]uint64, n)
   for i := range bits {
       bits[i] = make([]uint64, g)
   }
   // process edges
   for _, e := range edges {
       u, v := e.u, e.v
       // check common neighbor
       for k := 0; k < g; k++ {
           if bits[u][k]&bits[v][k] != 0 {
               r := math.Sqrt(float64(e.d)) / 2.0
               fmt.Printf("%.10f\n", r)
               return
           }
       }
       // add edge u-v
       bits[u][v>>6] |= 1 << (uint(v) & 63)
       bits[v][u>>6] |= 1 << (uint(u) & 63)
   }
   // no triple found
   fmt.Printf("0.0\n")
}
