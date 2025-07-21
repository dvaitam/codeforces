package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   var a, b int
   fmt.Fscan(in, &n, &m, &a, &b)
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ys[i])
   }
   yps := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(in, &yps[j])
   }
   ls := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(in, &ls[j])
   }

   // Precompute distance from O to Ai
   D := make([]float64, n)
   aa := float64(a)
   for i := 0; i < n; i++ {
       yi := float64(ys[i])
       D[i] = math.Hypot(aa, yi)
   }
   deltaX := float64(b - a)

   bestI := 0
   bestJ := 0
   minCost := math.Inf(1)

   // For each eastern bank point, find best western bank partner
   for j := 0; j < m; j++ {
       yj := float64(yps[j])
       // Move bestI while next gives better cost
       for bestI+1 < n {
           // cost at bestI
           dy0 := float64(ys[bestI]) - yj
           cost0 := D[bestI] + math.Hypot(deltaX, dy0)
           // cost at bestI+1
           dy1 := float64(ys[bestI+1]) - yj
           cost1 := D[bestI+1] + math.Hypot(deltaX, dy1)
           if cost1 <= cost0 {
               bestI++
           } else {
               break
           }
       }
       // compute current total cost including ls[j]
       dy := float64(ys[bestI]) - yj
       bridge := math.Hypot(deltaX, dy)
       total := D[bestI] + bridge + float64(ls[j])
       if total < minCost {
           minCost = total
           bestJ = j
       }
   }
   // Output 1-based indices: bestI+1 and bestJ+1
   fmt.Fprintf(out, "%d %d\n", bestI+1, bestJ+1)
}
