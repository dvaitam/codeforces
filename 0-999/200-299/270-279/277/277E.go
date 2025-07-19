package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   x := make([]int, n)
   y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }
   m := 2 * n
   const Inf = 1e18
   const Eps = 1e-6
   // weight matrix
   g := make([][]float64, n)
   for i := 0; i < n; i++ {
       g[i] = make([]float64, m)
       for j := 0; j < m; j++ {
           g[i][j] = -Inf
       }
   }
   // labels
   lx := make([]float64, n)
   ly := make([]float64, m)
   // slack
   slk := make([]float64, m)
   // matches
   match := make([]int, m)
   for i := range match {
       match[i] = -1
   }
   visx := make([]bool, n)
   visy := make([]bool, m)
   flw := 0

   // squared distance helper
   sqr := func(a int) float64 {
       return float64(a * a)
   }
   // distance (negative Euclidean)
   dist := func(i, j int) float64 {
       dx := x[i] - x[j]
       dy := y[i] - y[j]
       return -math.Sqrt(sqr(dx) + sqr(dy))
   }
   // add edge from i to j and j+n
   addEdge := func(i, j int, d float64) {
       g[i][j] = d
       g[i][j+n] = d
       if d > lx[i] {
           lx[i] = d
       }
   }

   // build graph
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if y[i] < y[j] {
               addEdge(i, j, dist(i, j))
           }
       }
   }

   // Hungarian DFS
   var hun func(u int) bool
   hun = func(u int) bool {
       visx[u] = true
       for v := 0; v < m; v++ {
           if g[u][v] <= -Inf/2 {
               continue
           }
           d := lx[u] + ly[v] - g[u][v]
           if d > Eps {
               if d < slk[v] {
                   slk[v] = d
               }
               continue
           }
           if visy[v] {
               continue
           }
           visy[v] = true
           if match[v] < 0 || hun(match[v]) {
               match[v] = u
               return true
           }
       }
       return false
   }

   // KM algorithm
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           slk[j] = Inf
       }
       for {
           for k := 0; k < n; k++ {
               visx[k] = false
           }
           for k := 0; k < m; k++ {
               visy[k] = false
           }
           if hun(i) {
               break
           }
           // adjust labels
           d := Inf
           for j := 0; j < m; j++ {
               if !visy[j] && slk[j] < d {
                   d = slk[j]
               }
           }
           if d >= Inf {
               flw++
               break
           }
           for j := 0; j < n; j++ {
               if visx[j] {
                   lx[j] -= d
               }
           }
           for j := 0; j < m; j++ {
               if visy[j] {
                   ly[j] += d
               }
           }
       }
       if flw > 1 {
           break
       }
   }
   if flw > 1 {
       fmt.Fprintln(writer, -1)
       return
   }
   var ans float64
   for v := 0; v < m; v++ {
       if match[v] >= 0 {
           ans -= g[match[v]][v]
       }
   }
   fmt.Fprintf(writer, "%.6f\n", ans)
}
