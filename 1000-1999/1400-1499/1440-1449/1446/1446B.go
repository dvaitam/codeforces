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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   var A, B string
   fmt.Fscan(reader, &A, &B)

   // dpPrev and dpCur are rolling arrays for local alignment dynamic programming
   dpPrev := make([]int, m+1)
   dpCur := make([]int, m+1)
   maxScore := 0
   for i := 1; i <= n; i++ {
       dpCur[0] = 0
       for j := 1; j <= m; j++ {
           // gap in A or B
           best := dpPrev[j] - 1
           if v := dpCur[j-1] - 1; v > best {
               best = v
           }
           // match gives +2
           if A[i-1] == B[j-1] {
               if v := dpPrev[j-1] + 2; v > best {
                   best = v
               }
           }
           if best < 0 {
               best = 0
           }
           dpCur[j] = best
           if best > maxScore {
               maxScore = best
           }
       }
       // swap previous and current rows
       dpPrev, dpCur = dpCur, dpPrev
   }

   fmt.Fprint(writer, maxScore)
}
