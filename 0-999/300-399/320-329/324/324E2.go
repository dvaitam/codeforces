package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   const mod = 1000000007
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   // adjacency: allowed edges where vision matches direct edge
   adj := make([][]int, n)
   for i := 0; i < n; i++ {
       adj[i] = make([]int, n)
   }
   for i := 0; i < m; i++ {
       var x, y, k int
       fmt.Fscan(reader, &x, &y, &k)
       v := make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &v[j])
       }
       // check if this edge produces vision exactly [x, y]
       if k == 2 && v[0] == x && v[1] == y {
           adj[x-1][y-1] = 1
       }
   }
   maxLen := 2 * n
   // dp[s][u]: number of walks of exactly s edges ending at u
   dp := make([][]int64, maxLen)
   for i := 0; i < maxLen; i++ {
       dp[i] = make([]int64, n)
   }
   // 0 edges: starting at any node
   for u := 0; u < n; u++ {
       dp[0][u] = 1
   }
   // build dp
   for s := 1; s < maxLen; s++ {
       for u := 0; u < n; u++ {
           var sum int64
           // previous v to u
           for v := 0; v < n; v++ {
               if adj[v][u] != 0 {
                   sum += dp[s-1][v]
                   if sum >= mod {
                       sum -= mod
                   }
               }
           }
           dp[s][u] = sum
       }
   }
   // output for lengths 1..2n
   for L := 1; L <= maxLen; L++ {
       if L == 1 {
           fmt.Fprintln(writer, 0)
           continue
       }
       s := L - 1
       var total int64
       for u := 0; u < n; u++ {
           total += dp[s][u]
           if total >= mod {
               total -= mod
           }
       }
       fmt.Fprintln(writer, total)
   }
}
