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
   const nax = 1 << 9
   fr := make([]int, nax)
   for i := 0; i < n; i++ {
       var k, x int
       fmt.Fscan(reader, &k)
       mask := 0
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &x)
           mask |= 1 << (x - 1)
       }
       fr[mask]++
   }
   ile := make([]int, nax)
   for i := 0; i < nax; i++ {
       if fr[i] == 0 {
           continue
       }
       for x := 0; x < nax; x++ {
           if x&i == i {
               ile[x] += fr[i]
           }
       }
   }
   const INF = int(1e9 + 5)
   type pair struct{ first, second int }
   piz := make([]pair, nax)
   for i := range piz {
       piz[i] = pair{INF, INF}
   }
   bon := pair{INF, INF}
   for idx := 1; idx <= m; idx++ {
       var pr, k, x int
       fmt.Fscan(reader, &pr, &k)
       mask := 0
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &x)
           mask |= 1 << (x - 1)
       }
       if pr < piz[mask].first {
           if piz[mask].first != INF {
               old := piz[mask]
               if old.first < bon.first || (old.first == bon.first && old.second < bon.second) {
                   bon = old
               }
           }
           piz[mask] = pair{pr, idx}
       } else {
           cur := pair{pr, idx}
           if cur.first < bon.first || (cur.first == bon.first && cur.second < bon.second) {
               bon = cur
           }
       }
   }
   ans := pair{INF, INF}
   var ret pair
   for i := 0; i < nax; i++ {
       if piz[i].first == INF {
           continue
       }
       // combine with different masks
       for j := 0; j < nax; j++ {
           if i == j || piz[j].first == INF {
               continue
           }
           cm := i | j
           cov := ile[cm]
           sum := piz[i].first + piz[j].first
           nw := pair{-cov, sum}
           if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
               ans = nw
               ret = pair{piz[i].second, piz[j].second}
           }
       }
       // combine with second best
       if bon.first != INF {
           cov := ile[i]
           sum := piz[i].first + bon.first
           nw := pair{-cov, sum}
           if nw.first < ans.first || (nw.first == ans.first && nw.second < ans.second) {
               ans = nw
               ret = pair{piz[i].second, bon.second}
           }
       }
   }
   fmt.Fprintln(writer, ret.first, ret.second)
}
