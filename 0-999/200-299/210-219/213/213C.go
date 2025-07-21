package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       a[i] = make([]int, n+1)
       for j := 1; j <= n; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   // dp for current and previous steps
   const INF = -1 << 60
   dpPrev := make([][]int, n+1)
   dpCur := make([][]int, n+1)
   for i := 0; i <= n; i++ {
       dpPrev[i] = make([]int, n+1)
       dpCur[i] = make([]int, n+1)
       for j := 0; j <= n; j++ {
           dpPrev[i][j] = INF
           dpCur[i][j] = INF
       }
   }
   // initial at step 2 (1+1)
   dpPrev[1][1] = a[1][1]
   // steps from 3 to 2n
   for s := 3; s <= 2*n; s++ {
       // reset dpCur
       for i := 1; i <= n; i++ {
           for j := 1; j <= n; j++ {
               dpCur[i][j] = INF
           }
       }
       // i1 and i2 are x-coordinates of two walkers
       for i1 := max(1, s-n); i1 <= n && i1 < s; i1++ {
           j1 := s - i1
           if j1 < 1 || j1 > n {
               continue
           }
           for i2 := max(1, s-n); i2 <= n && i2 < s; i2++ {
               j2 := s - i2
               if j2 < 1 || j2 > n {
                   continue
               }
               best := dpPrev[i1][i2] // both came from left
               best = max(best, dpPrev[i1-1][i2])   // down, left
               best = max(best, dpPrev[i1][i2-1])   // left, down
               best = max(best, dpPrev[i1-1][i2-1]) // both down
               if best <= INF/2 {
                   continue
               }
               val := best + a[i1][j1]
               if i1 != i2 || j1 != j2 {
                   val += a[i2][j2]
               }
               dpCur[i1][i2] = val
           }
       }
       // swap
       dpPrev, dpCur = dpCur, dpPrev
   }
   // result at both at (n,n), step 2n
   res := dpPrev[n][n]
   fmt.Fprintln(writer, res)
}
