package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct{
   v, p int
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, k, w int
   fmt.Fscan(in, &n, &m, &k, &w)
   // read levels
   field := make([][]string, k)
   for i := 0; i < k; i++ {
       field[i] = make([]string, n)
       for r := 0; r < n; r++ {
           fmt.Fscan(in, &field[i][r])
       }
   }
   // compute pairwise distances
   maxCost := n * m
   dist := make([][]int, k)
   for i := 0; i < k; i++ {
       dist[i] = make([]int, k)
   }
   for i := 0; i < k; i++ {
       for j := i + 1; j < k; j++ {
           diff := 0
           for r := 0; r < n; r++ {
               a := field[i][r]
               b := field[j][r]
               for c := 0; c < m; c++ {
                   if a[c] != b[c] {
                       diff++
                   }
               }
           }
           cost := min(maxCost, diff * w)
           dist[i][j] = cost
           dist[j][i] = cost
       }
   }
   // Prim's algorithm
   used := make([]bool, k)
   minCost := make([]int, k)
   prev := make([]int, k)
   // start from 0
   used[0] = true
   total := maxCost
   ans := make([]pair, 0, k)
   ans = append(ans, pair{0, -1})
   // init costs
   for i := 1; i < k; i++ {
       minCost[i] = dist[0][i]
       prev[i] = 0
   }
   // build MST
   for cnt := 1; cnt < k; cnt++ {
       v := -1
       for i := 0; i < k; i++ {
           if !used[i] && (v == -1 || minCost[i] < minCost[v]) {
               v = i
           }
       }
       if minCost[v] >= maxCost {
           total += maxCost
           ans = append(ans, pair{v, -1})
       } else {
           total += minCost[v]
           ans = append(ans, pair{v, prev[v]})
       }
       used[v] = true
       // update neighbors
       for u := 0; u < k; u++ {
           if !used[u] && dist[v][u] < minCost[u] {
               minCost[u] = dist[v][u]
               prev[u] = v
           }
       }
   }
   // output
   fmt.Fprintln(out, total)
   for _, pr := range ans {
       // print 1-based indices; parent -1 -> 0
       x := pr.v + 1
       y := pr.p + 1
       if pr.p < 0 {
           y = 0
       }
       fmt.Fprintln(out, x, y)
   }
}
