package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M, si, sj int
   if _, err := fmt.Fscan(in, &N, &M); err != nil {
       return
   }
   fmt.Fscan(in, &si, &sj)
   // zero-based indices
   si--
   sj--
   // if already at bottom row
   if si == N-1 {
       fmt.Printf("0.0000\n")
       return
   }
   // D holds expected steps for next row (row r+1)
   D := make([]float64, M)
   // temporary arrays for Thomas algorithm
   a := make([]float64, M)
   b := make([]float64, M)
   c := make([]float64, M)
   d := make([]float64, M)
   cp := make([]float64, M)
   dp := make([]float64, M)
   x := make([]float64, M)
   // process rows from bottom-1 down to 0
   for r := N - 2; r >= 0; r-- {
       // build tridiagonal system for row r
       for j := 0; j < M; j++ {
           down := D[j]
           switch {
           case j == 0 && M > 1:
               a[j] = 0
               b[j] = 2
               c[j] = -1
               d[j] = 3 + down
           case j == M-1 && M > 1:
               a[j] = -1
               b[j] = 2
               c[j] = 0
               d[j] = 3 + down
           case M == 1:
               // only one column
               a[j] = 0
               b[j] = 1    // deg-1 = 1
               c[j] = 0
               // deg=2
               d[j] = 2 + down
           default:
               // middle columns
               a[j] = -1
               b[j] = 3
               c[j] = -1
               d[j] = 4 + down
           }
       }
       // forward sweep
       cp[0] = c[0] / b[0]
       dp[0] = d[0] / b[0]
       for j := 1; j < M; j++ {
           denom := b[j] - a[j]*cp[j-1]
           cp[j] = c[j] / denom
           dp[j] = (d[j] - a[j]*dp[j-1]) / denom
       }
       // back substitution
       x[M-1] = dp[M-1]
       for j := M - 2; j >= 0; j-- {
           x[j] = dp[j] - cp[j]*x[j+1]
       }
       // move x to D
       copy(D, x)
   }
   // result at D[sj]
   ans := D[sj]
   fmt.Printf("%.4f\n", ans)
}
