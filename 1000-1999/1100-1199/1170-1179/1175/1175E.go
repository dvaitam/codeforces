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
   const MAXC = 500000
   best := make([]int, MAXC+2)
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       if best[l] < r {
           best[l] = r
       }
   }
   // prefix max to get farthest reach from each position
   for i := 1; i <= MAXC; i++ {
       if best[i] < best[i-1] {
           best[i] = best[i-1]
       }
   }
   // build binary lifting table
   const LOG = 19
   // nxt[k][i]: farthest reachable from i using 2^k intervals
   nxt := make([][]int, LOG)
   nxt[0] = make([]int, MAXC+2)
   copy(nxt[0], best)
   for k := 1; k < LOG; k++ {
       nxt[k] = make([]int, MAXC+2)
       for i := 0; i <= MAXC; i++ {
           nxt[k][i] = nxt[k-1][ nxt[k-1][i] ]
       }
   }
   // answer queries
   for qi := 0; qi < m; qi++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       // greedy binary lifting
       curr := x
       ans := 0
       if curr < 0 {
           curr = 0
       }
       for k := LOG - 1; k >= 0; k-- {
           if nxt[k][curr] < y {
               ans += 1 << k
               curr = nxt[k][curr]
           }
       }
       // one more interval
       if best[curr] >= y {
           ans++
           fmt.Fprint(writer, ans)
       } else {
           fmt.Fprint(writer, -1)
       }
       if qi+1 < m {
           fmt.Fprint(writer, " ")
       }
   }
   writer.WriteByte('\n')
}
