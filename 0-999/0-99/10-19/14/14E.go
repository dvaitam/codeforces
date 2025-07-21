package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, t int
   if _, err := fmt.Fscan(reader, &n, &t); err != nil {
       return
   }
   // dp[pos][lastY][peaks][valleys][lastDir]
   // lastDir: 0 undefined, 1 up, 2 down
   // maximum n=20, t<=10
   valleyTarget := t - 1
   // dp dimensions: pos from 1..n
   dp := make([][][][][]int64, n+1)
   for i := range dp {
       dp[i] = make([][][][]int64, 5)
       for y := 0; y < 5; y++ {
           dp[i][y] = make([][][]int64, t+1)
           for p := 0; p <= t; p++ {
               dp[i][y][p] = make([][]int64, t+1)
               for v := 0; v <= t; v++ {
                   dp[i][y][p][v] = make([]int64, 3)
               }
           }
       }
   }
   // init at pos 1
   for y := 1; y <= 4; y++ {
       dp[1][y][0][0][0] = 1
   }
   for pos := 1; pos < n; pos++ {
       for lastY := 1; lastY <= 4; lastY++ {
           for p := 0; p <= t; p++ {
               for v := 0; v <= valleyTarget; v++ {
                   for lastDir := 0; lastDir < 3; lastDir++ {
                       cnt := dp[pos][lastY][p][v][lastDir]
                       if cnt == 0 {
                           continue
                       }
                       // choose next y
                       for y2 := 1; y2 <= 4; y2++ {
                           if y2 == lastY {
                               continue
                           }
                           // direction
                           var dir int
                           if y2 > lastY {
                               dir = 1
                           } else {
                               dir = 2
                           }
                           np, nv := p, v
                           if lastDir == 1 && dir == 2 {
                               np++
                           } else if lastDir == 2 && dir == 1 {
                               nv++
                           }
                           if np > t || nv > valleyTarget {
                               continue
                           }
                           dp[pos+1][y2][np][nv][dir] += cnt
                       }
                   }
               }
           }
       }
   }
   var res int64
   for y := 1; y <= 4; y++ {
       for lastDir := 1; lastDir <= 2; lastDir++ {
           res += dp[n][y][t][valleyTarget][lastDir]
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, res)
}
