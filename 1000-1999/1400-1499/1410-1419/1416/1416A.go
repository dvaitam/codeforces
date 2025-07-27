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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // last position of value, max gap between occurrences
       lastPos := make([]int, n+1)
       maxGap := make([]int, n+1)
       seen := make([]bool, n+1)
       for i := 1; i <= n; i++ {
           x := a[i]
           gap := i - lastPos[x]
           if gap > maxGap[x] {
               maxGap[x] = gap
           }
           lastPos[x] = i
           seen[x] = true
       }
       for x := 1; x <= n; x++ {
           if !seen[x] {
               continue
           }
           // tail gap to end
           gap := (n + 1) - lastPos[x]
           if gap > maxGap[x] {
               maxGap[x] = gap
           }
       }
       const inf = int(1e9)
       best := make([]int, n+2)
       for i := 1; i <= n+1; i++ {
           best[i] = inf
       }
       for x := 1; x <= n; x++ {
           if !seen[x] {
               continue
           }
           k := maxGap[x]
           if x < best[k] {
               best[k] = x
           }
       }
       // build answers: prefix minima
       res := make([]int, n+1)
       curr := inf
       for k := 1; k <= n; k++ {
           if best[k] < curr {
               curr = best[k]
           }
           if curr == inf {
               res[k] = -1
           } else {
               res[k] = curr
           }
       }
       // output
       for k := 1; k <= n; k++ {
           if k > 1 {
               writer.WriteByte(' ')
           }
           fmt.Fprintf(writer, "%d", res[k])
       }
       writer.WriteByte('\n')
   }
}
