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
   var n, k int
   fmt.Fscan(reader, &n, &k)
   l := make([]uint64, n)
   r := make([]uint64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &l[i], &r[i])
   }
   c := make([]int64, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &c[i])
   }
   var res int64
   const INF = int64(1e18)
   // process each bit independently
   for j := 0; j < k; j++ {
       // dpPrev[v] = min transitions up to previous position ending with bit v
       dpPrev := [2]int64{INF, INF}
       // initial position
       bL := (l[0] >> j) & 1
       bR := (r[0] >> j) & 1
       if bL == bR {
           dpPrev[bL] = 0
       } else {
           dpPrev[0], dpPrev[1] = 0, 0
       }
       // dp over positions
       for i := 1; i < n; i++ {
           dpCurr := [2]int64{INF, INF}
           bL = (l[i] >> j) & 1
           bR = (r[i] >> j) & 1
           for v := 0; v < 2; v++ {
               // check if bit v is allowed at position i
               if bL == bR && v != int(bL) {
                   continue
               }
               for u := 0; u < 2; u++ {
                   if dpPrev[u] == INF {
                       continue
                   }
                   cost := dpPrev[u]
                   if u != v {
                       cost++
                   }
                   if cost < dpCurr[v] {
                       dpCurr[v] = cost
                   }
               }
           }
           dpPrev = dpCurr
       }
       // pick best ending bit
       best := dpPrev[0]
       if dpPrev[1] < best {
           best = dpPrev[1]
       }
       res += best * c[j]
   }
   fmt.Fprint(writer, res)
}
