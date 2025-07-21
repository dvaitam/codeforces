package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   l := make([]int, n+1)
   r := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &l[i], &r[i])
   }
   // dp[i][j]: openings from i to j can form a valid sequence
   dp := make([][]bool, n+2)
   choice := make([][]int, n+2)
   for i := 0; i <= n+1; i++ {
       dp[i] = make([]bool, n+2)
       choice[i] = make([]int, n+2)
   }
   // empty segments
   for i := 1; i <= n+1; i++ {
       dp[i][i-1] = true
   }
   // build for lengths
   for length := 1; length <= n; length++ {
       for i := 1; i+length-1 <= n; i++ {
           j := i + length - 1
           // try matching opening i with some partner at i+t
           // t = number of inner openings
           tmin := l[i] / 2
           tmax := (r[i] - 1) / 2
           if tmin < 0 {
               tmin = 0
           }
           if tmax > j-i {
               tmax = j - i
           }
           for t := tmin; t <= tmax; t++ {
               // check inside and rest
               if dp[i+1][i+t] && dp[i+t+1][j] {
                   dp[i][j] = true
                   choice[i][j] = t
                   break
               }
           }
       }
   }
   if !dp[1][n] {
       fmt.Println("IMPOSSIBLE")
       return
   }
   // reconstruct
   buf := make([]byte, 0, 2*n)
   var build func(i, j int)
   build = func(i, j int) {
       if i > j {
           return
       }
       t := choice[i][j]
       buf = append(buf, '(')
       build(i+1, i+t)
       buf = append(buf, ')')
       build(i+t+1, j)
   }
   build(1, n)
   fmt.Println(string(buf))
}
