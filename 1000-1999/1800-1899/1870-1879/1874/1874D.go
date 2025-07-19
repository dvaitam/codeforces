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
   // dp array: f[i][j] = min cost to reach row i at column j
   const INF = 1e18
   f := make([][]float64, n+1)
   for i := 1; i <= n; i++ {
       f[i] = make([]float64, m+1)
       for j := 1; j <= m; j++ {
           f[i][j] = INF
       }
   }
   // base case
   for j := 1; j <= m; j++ {
       f[1][j] = 0
   }
   // dp transitions
   for i := 1; i < n; i++ {
       step := n - i
       for j := 1; j <= m; j++ {
           val := f[i][j]
           if val >= INF {
               continue
           }
           maxK := (m - j) / step
           for k := 1; k <= maxK; k++ {
               nj := j + k*step
               cost := val + float64(j)/float64(k)
               if cost < f[i+1][nj] {
                   f[i+1][nj] = cost
               }
           }
       }
   }
   // find answer
   ans := INF
   for j := 1; j <= m; j++ {
       if f[n][j] < ans {
           ans = f[n][j]
       }
   }
   res := float64(n) + 2*ans
   fmt.Printf("%.12f\n", res)
}
