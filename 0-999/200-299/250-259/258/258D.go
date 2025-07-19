package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   w := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &w[i])
   }
   // P[i][j]: probability that pair (i,j) is inverted, 0-based
   P := make([][]float64, n)
   for i := 0; i < n; i++ {
       P[i] = make([]float64, n)
   }
   // initial probabilities
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if w[i] > w[j] {
               P[i][j] = 1.0
           }
       }
   }
   for k := 0; k < m; k++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       a--
       b--
       if a > b {
           a, b = b, a
       }
       // after this move, positions a and b may be swapped with prob 0.5
       P[a][b] = 0.5
       // update pairs
       for j := 0; j < a; j++ {
           P[j][a] = (P[j][a] + P[j][b]) * 0.5
           P[j][b] = P[j][a]
       }
       for j := b + 1; j < n; j++ {
           P[a][j] = (P[a][j] + P[b][j]) * 0.5
           P[b][j] = P[a][j]
       }
       for j := a + 1; j < b; j++ {
           t1 := P[a][j]
           t2 := P[j][b]
           P[a][j] = (1.0 - t2 + t1) * 0.5
           P[j][b] = (t2 + 1.0 - t1) * 0.5
       }
   }
   var R float64
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           R += P[i][j]
       }
   }
   fmt.Printf("%.6f\n", R)
}
