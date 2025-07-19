package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k, x int
   fmt.Fscan(reader, &n, &k, &x)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   const INF = int64(-1e18)
   // d[i][j]: max sum using j picks ending at position i
   d := make([][]int64, n+1)
   for i := 0; i <= n; i++ {
       d[i] = make([]int64, x+1)
       for j := 0; j <= x; j++ {
           d[i][j] = INF
       }
   }
   d[0][0] = 0

   for i := 0; i < n; i++ {
       for j := 0; j < x; j++ {
           if d[i][j] != INF {
               for z := i + 1; z <= i+k && z <= n; z++ {
                   val := d[i][j] + a[z]
                   if val > d[z][j+1] {
                       d[z][j+1] = val
                   }
               }
           }
       }
   }

   ans := INF
   // answer is max over last k positions
   for i := n - k + 1; i <= n; i++ {
       if i >= 0 && i <= n && d[i][x] > ans {
           ans = d[i][x]
       }
   }
   if ans == INF {
       fmt.Fprint(writer, -1)
   } else {
       fmt.Fprint(writer, ans)
   }
}
