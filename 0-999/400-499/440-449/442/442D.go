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
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   parent := make([]int, n+2)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &parent[i+1])
   }
   // initialize arrays for nodes 1..n+1
   N := n + 1
   size := make([]int, N+1)
   dp := make([]int, N+1)
   bestNon := make([]int, N+1)
   heavy1 := make([]int, N+1)
   heavy2 := make([]int, N+1)
   // default bestNon = -1 (no non-heavy child)
   for i := 1; i <= N; i++ {
       bestNon[i] = -1
   }
   // root is 1, parent[1] == 0
   // process additions
   for i := 1; i <= n; i++ {
       v := i + 1
       size[v] = 1
       dp[v] = 0
       // bestNon[v] already -1, heavy1/2[v] = 0
       prev := v
       u := parent[v]
       // update up to root
       for u != 0 {
           size[u]++
           heavyChanged := false
           // capacity: 2 for root, else 1
           if u == 1 {
               // root: up to two heavy children by size
               if heavy1[u] == 0 {
                   heavy1[u] = prev
                   heavyChanged = true
               } else if heavy2[u] == 0 {
                   // pick order by size
                   if size[prev] > size[heavy1[u]] {
                       heavy2[u] = heavy1[u]
                       heavy1[u] = prev
                   } else {
                       heavy2[u] = prev
                   }
                   heavyChanged = true
               } else {
                   // find smallest heavy
                   var minh int
                   if size[heavy1[u]] < size[heavy2[u]] {
                       minh = heavy1[u]
                   } else {
                       minh = heavy2[u]
                   }
                   if size[prev] > size[minh] {
                       // replace smaller heavy
                       if minh == heavy1[u] {
                           heavy1[u] = prev
                       } else {
                           heavy2[u] = prev
                       }
                       heavyChanged = true
                       // old heavy becomes non-heavy
                       bestNon[u] = max(bestNon[u], dp[minh])
                   } else {
                       // prev remains non-heavy
                       bestNon[u] = max(bestNon[u], dp[prev])
                   }
               }
           } else {
               // non-root: at most one heavy child
               if heavy1[u] == 0 {
                   heavy1[u] = prev
                   heavyChanged = true
               } else if size[prev] > size[heavy1[u]] {
                   // swap heavy
                   old := heavy1[u]
                   heavy1[u] = prev
                   heavyChanged = true
                   // old heavy becomes non-heavy
                   bestNon[u] = max(bestNon[u], dp[old])
               } else {
                   // non-heavy
                   bestNon[u] = max(bestNon[u], dp[prev])
               }
           }
           // recompute dp[u]
           oldDp := dp[u]
           // best heavy dp
           bestH := -1
           if heavy1[u] != 0 {
               bestH = max(bestH, dp[heavy1[u]])
           }
           if u == 1 && heavy2[u] != 0 {
               bestH = max(bestH, dp[heavy2[u]])
           }
           // compute new dp
           ndp := bestH
           if u == 1 {
               // root: non-heavy edges do not count as light
               ndp = max(ndp, bestNon[u])
           } else {
               if bestNon[u] >= 0 {
                   ndp = max(ndp, bestNon[u]+1)
               }
           }
           if ndp < 0 {
               ndp = 0
           }
           dp[u] = ndp
           // stop if nothing changed
           if dp[u] == oldDp && !heavyChanged {
               break
           }
           prev = u
           u = parent[u]
       }
       // output cost = dp[1] + 1
       fmt.Fprint(out, dp[1]+1)
       if i < n {
           out.WriteByte(' ')
       }
   }
   out.WriteByte('\n')
}
