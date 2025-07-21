package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, t int
   if _, err := fmt.Fscan(reader, &n, &m, &t); err != nil {
       return
   }
   // Precompute combinations C[i][j] for 0 <= i,j <= max(n,m)
   maxNM := n
   if m > maxNM {
       maxNM = m
   }
   C := make([][]int64, maxNM+1)
   for i := 0; i <= maxNM; i++ {
       C[i] = make([]int64, i+1)
       C[i][0] = 1
       for j := 1; j <= i; j++ {
           if j == i {
               C[i][j] = 1
           } else {
               C[i][j] = C[i-1][j-1] + C[i-1][j]
           }
       }
   }
   var result int64 = 0
   // choose b boys and g girls, where b+g = t, b>=4, g>=1
   for b := 4; b <= n; b++ {
       g := t - b
       if g < 1 || g > m {
           continue
       }
       // C(n, b) * C(m, g)
       result += C[n][b] * C[m][g]
   }
   fmt.Println(result)
}
