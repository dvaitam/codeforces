package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var ax, ay, bx, by, tx, ty float64
   var n int
   if _, err := fmt.Fscan(reader, &ax, &ay, &bx, &by, &tx, &ty, &n); err != nil {
       return
   }
   x := make([]float64, n)
   y := make([]float64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }
   // total round-trip distance from trash (tx,ty) to all points
   sum := 0.0
   distT := make([]float64, n)
   for i := 0; i < n; i++ {
       d := math.Hypot(x[i]-tx, y[i]-ty)
       distT[i] = d
       sum += 2 * d
   }
   // base answer: one person picks exactly one, the other picks at most one
   ans := math.Inf(1)
   // consider only one special pick (if n==1, both loops handle)
   for i := 0; i < n; i++ {
       // A picks i
       costA := math.Hypot(x[i]-ax, y[i]-ay)
       cand := costA + sum - distT[i]
       if cand < ans {
           ans = cand
       }
       // B picks i
       costB := math.Hypot(x[i]-bx, y[i]-by)
       cand = costB + sum - distT[i]
       if cand < ans {
           ans = cand
       }
   }
   if n == 1 {
       fmt.Printf("%.10f\n", ans)
       return
   }
   // deltas for two special picks
   deltaA := make([]float64, n)
   deltaB := make([]float64, n)
   for i := 0; i < n; i++ {
       deltaA[i] = math.Hypot(x[i]-ax, y[i]-ay) - distT[i]
       deltaB[i] = math.Hypot(x[i]-bx, y[i]-by) - distT[i]
   }
   // sort indices by deltaB ascending
   idx := make([]int, n)
   for i := range idx {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       return deltaB[idx[i]] < deltaB[idx[j]]
   })
   first := idx[0]
   second := idx[1]
   // try combinations: A picks i, B picks best or second best
   for i := 0; i < n; i++ {
       total := sum + deltaA[i]
       if i == first {
           total += deltaB[second]
       } else {
           total += deltaB[first]
       }
       if total < ans {
           ans = total
       }
   }
   fmt.Printf("%.10f\n", ans)
}
