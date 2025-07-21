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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 2; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Nodes for pos=2..n, dir=0,1: id = (pos-2)*2 + dir
   m := (n - 1) * 2
   const UNVIS = -2
   const INF = -1
   dp := make([]int64, m)
   rev := make([][]int, m)
   for i := 0; i < m; i++ {
       dp[i] = UNVIS
   }
   // queue
   q := make([]int, 0, m)
   // build graph
   for pos := 2; pos <= n; pos++ {
       idx := (pos - 2) * 2
       for dir := 0; dir < 2; dir++ {
           u := idx + dir
           var vpos int64
           if dir == 0 {
               vpos = int64(pos) + a[pos]
           } else {
               vpos = int64(pos) - a[pos]
           }
           vdir := dir ^ 1
           // if vpos in [2..n], add reverse edge
           if vpos >= 2 && vpos <= int64(n) {
               vidx := (int(vpos) - 2) * 2
               v := vidx + vdir
               rev[v] = append(rev[v], u)
           } else if vpos <= 0 || vpos > int64(n) {
               // terminal exit
               dp[u] = a[pos]
               q = append(q, u)
           }
           // else vpos == 1: returns to start, treated as cycle -> skip
       }
   }
   // BFS propagate
   for head := 0; head < len(q); head++ {
       v := q[head]
       for _, u := range rev[v] {
           if dp[u] != UNVIS {
               continue
           }
           // from u, weight is a[pos(u)]
           pos := (u/2 + 2)
           if dp[v] == INF {
               dp[u] = INF
           } else {
               dp[u] = a[pos] + dp[v]
           }
           q = append(q, u)
       }
   }
   // remaining unvisited are cycles => INF
   for i := 0; i < m; i++ {
       if dp[i] == UNVIS {
           dp[i] = INF
       }
   }
   // output answers for i = 1..n-1
   for i := 1; i <= n-1; i++ {
       pos1 := 1 + i
       // dir=1
       u := (pos1-2)*2 + 1
       if dp[u] == INF {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, dp[u] + int64(i))
       }
   }
}
