package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var T int
   fmt.Fscan(in, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       ans := make([]int, 0, n*(n-1)/2)
       if n&1 == 1 {
           half := n / 2
           for i := 1; i <= n; i++ {
               for j := i + 1; j <= n; j++ {
                   if j-i <= half {
                       ans = append(ans, 1)
                   } else {
                       ans = append(ans, -1)
                   }
               }
           }
       } else {
           for i := 1; i < n; i += 2 {
               // match between i and i+1 is a tie
               ans = append(ans, 0)
               m := n - i - 1
               half := m / 2
               // results for matches (i, j)
               for k := 1; k <= m; k++ {
                   if k <= half {
                       ans = append(ans, 1)
                   } else {
                       ans = append(ans, -1)
                   }
               }
               // results for matches (i+1, j)
               for k := 1; k <= m; k++ {
                   if k <= half {
                       ans = append(ans, -1)
                   } else {
                       ans = append(ans, 1)
                   }
               }
           }
       }
       // output results
       for idx, v := range ans {
           if idx > 0 {
               out.WriteByte(' ')
           }
           switch v {
           case 1:
               out.WriteString("1")
           case 0:
               out.WriteString("0")
           case -1:
               out.WriteString("-1")
           }
       }
       out.WriteByte('\n')
   }
}
