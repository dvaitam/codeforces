package main

import (
   "bufio"
   "fmt"
   "os"
)

func minFloat(a, b float64) float64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   c := make([]int, n)
   total := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &c[i])
       total += c[i]
   }

   // Binomial coefficients
   C := make([][]float64, n+1)
   for i := 0; i <= n; i++ {
       C[i] = make([]float64, i+1)
       C[i][0], C[i][i] = 1.0, 1.0
       for j := 1; j < i; j++ {
           C[i][j] = C[i-1][j] + C[i-1][j-1]
       }
   }

   // dp[j][p]: number of ways to pick j items with sum p
   dp := make([][]float64, n+1)
   for i := 0; i <= n; i++ {
       dp[i] = make([]float64, total+1)
   }
   dp[0][0] = 1.0
   now := 0
   for i := 0; i < n; i++ {
       ci := c[i]
       now += ci
       // update dp in reverse
       for j := i + 1; j >= 1; j-- {
           src := dp[j-1]
           dst := dp[j]
           for p := now; p >= ci; p-- {
               dst[p] += src[p-ci]
           }
       }
   }

   var ans float64
   nf := float64(n)
   xf := float64(x)
   for i := 1; i <= n; i++ {
       invCi := 1.0 / C[n][i]
       idx := float64(i)
       limit := (nf/idx + 1.0) * xf / 2.0
       for j := 1; j <= now; j++ {
           ways := dp[i][j]
           if ways == 0 {
               continue
           }
           value := float64(j) / idx
           cost := minFloat(limit, value)
           ans += ways * invCi * cost
       }
   }

   fmt.Fprintf(writer, "%.10f\n", ans)
}
