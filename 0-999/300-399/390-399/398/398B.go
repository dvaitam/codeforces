package main

import (
   "bufio"
   "fmt"
   "os"
)

const N = 2005

var (
   f [N][N]float64
   a, b [N]bool
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       if x >= 1 && x < N {
           a[x] = true
       }
       if y >= 1 && y < N {
           b[y] = true
       }
   }
   A, B := 0, 0
   for i := 1; i <= n; i++ {
       if a[i] {
           A++
       }
       if b[i] {
           B++
       }
   }
   nn := float64(n)
   nn2 := nn * nn
   // DP from end states
   for i := n; i >= A; i-- {
       for j := n; j >= B; j-- {
           if i < n || j < n {
               // number of pairs remaining: n*n - i*j
               denom := nn2 - float64(i*j)
               // expected value recurrence
               term1 := float64(n-i) * float64(j) * f[i+1][j]
               term2 := float64(i) * float64(n-j) * f[i][j+1]
               term3 := float64(n-i) * float64(n-j) * f[i+1][j+1]
               f[i][j] = (term1 + term2 + term3 + nn2) / denom
           }
       }
   }
   fmt.Printf("%.9f\n", f[A][B])
}
