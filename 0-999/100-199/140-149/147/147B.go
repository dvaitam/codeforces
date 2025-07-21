package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   const INF = int64(1e18)
   type edge struct{u, v int; w int64}
   edges := make([]edge, 0, m*2)
   for i := 0; i < m; i++ {
       var a, b int
       var cab, cba int
       fmt.Fscan(reader, &a, &b, &cab, &cba)
       a--
       b--
       edges = append(edges, edge{a, b, int64(cab)})
       edges = append(edges, edge{b, a, int64(cba)})
   }
   // dpPrev[i][j] = max weight of path from i to j with current length
   dpPrev := make([][]int64, n)
   dpCur := make([][]int64, n)
   for i := 0; i < n; i++ {
       dpPrev[i] = make([]int64, n)
       dpCur[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           dpPrev[i][j] = -INF
           dpCur[i][j] = -INF
       }
   }
   // initialize length-1 paths
   for _, e := range edges {
       if e.w > dpPrev[e.u][e.v] {
           dpPrev[e.u][e.v] = e.w
       }
   }
   // check cycles of length >=2
   for k := 2; k <= n; k++ {
       // reset dpCur
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               dpCur[i][j] = -INF
           }
       }
       // relax paths
       for _, e := range edges {
           u, v, w := e.u, e.v, e.w
           for i := 0; i < n; i++ {
               if dpPrev[i][u] > -INF {
                   val := dpPrev[i][u] + w
                   if val > dpCur[i][v] {
                       dpCur[i][v] = val
                   }
               }
           }
       }
       // check for positive cycle of length k
       for i := 0; i < n; i++ {
           if dpCur[i][i] > 0 {
               fmt.Println(k)
               return
           }
       }
       // prepare for next iteration
       dpPrev, dpCur = dpCur, dpPrev
   }
   // no positive cycle
   fmt.Println(0)
}
