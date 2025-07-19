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
   a := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Precompute powers p[k] = (1/m)^k
   p := make([]float64, n+1)
   p[0] = 1.0
   invM := 1.0 / float64(m)
   for i := 1; i <= n; i++ {
       p[i] = p[i-1] * invM
   }
   // Precompute binomial coefficients c[i][j]
   c := make([][]float64, n+1)
   for i := 0; i <= n; i++ {
       c[i] = make([]float64, n+1)
       c[i][0] = 1.0
       for j := 1; j <= i; j++ {
           c[i][j] = c[i-1][j-1] + c[i-1][j]
       }
   }
   // f[i][j]: probability using first i bins to allocate j balls
   f := make([][]float64, m+1)
   for i := 0; i <= m; i++ {
       f[i] = make([]float64, n+1)
   }
   var ans, pre float64
   // kk: maximum allowed rounds
   for kk := 1; kk <= n; kk++ {
       // reset f
       for i := 0; i <= m; i++ {
           for j := 0; j <= n; j++ {
               f[i][j] = 0.0
           }
       }
       f[0][0] = 1.0
       // DP over bins
       for i := 0; i < m; i++ {
           for j := 0; j <= n; j++ {
               if f[i][j] == 0.0 {
                   continue
               }
               maxk := a[i] * kk
               if maxk > n-j {
                   maxk = n - j
               }
               for k := 0; k <= maxk; k++ {
                   f[i+1][j+k] += f[i][j] * p[k] * c[n-j][k]
               }
           }
       }
       cur := f[m][n]
       ans += (cur - pre) * float64(kk)
       pre = cur
   }
   // Output with high precision
   fmt.Printf("%.20f\n", ans)
}
