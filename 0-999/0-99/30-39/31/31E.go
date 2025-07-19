package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(in, &n, &s); err != nil {
       return
   }
   // digits of A
   total := 2 * n
   digits := make([]int, total)
   for i := 0; i < total; i++ {
       digits[i] = int(s[i] - '0')
   }
   // powers of 10
   p10 := make([]int64, n+1)
   p10[0] = 1
   for i := 1; i <= n; i++ {
       p10[i] = p10[i-1] * 10
   }
   // dp f[i][j] = max sum after i moves by H and j by M
   f := make([][]int64, n+1)
   for i := range f {
       f[i] = make([]int64, n+1)
   }
   // transitions
   for i := 0; i <= n; i++ {
       for j := 0; j <= n; j++ {
           if i < n {
               idx := i + j
               // Homer picks
               w := int64(digits[idx]) * p10[n-i-1]
               f[i+1][j] = max(f[i+1][j], f[i][j] + w)
           }
           if j < n {
               idx := i + j
               // Marge picks
               w := int64(digits[idx]) * p10[n-j-1]
               f[i][j+1] = max(f[i][j+1], f[i][j] + w)
           }
       }
   }
   // reconstruct moves
   res := make([]byte, total)
   x, y := n, n
   for x > 0 || y > 0 {
       idx := x + y - 1
       if x > 0 {
           w := int64(digits[idx]) * p10[n-x]
           if f[x][y] == f[x-1][y] + w {
               res[idx] = 'H'
               x--
               continue
           }
       }
       // otherwise Marge
       res[idx] = 'M'
       y--
   }
   fmt.Println(string(res))
}
