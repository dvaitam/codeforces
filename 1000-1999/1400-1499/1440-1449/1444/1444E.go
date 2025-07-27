package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   in  = bufio.NewReader(os.Stdin)
   out = bufio.NewWriter(os.Stdout)
)

func flush() {
   out.Flush()
}

func main() {
   defer flush()
   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   alive := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       alive[i] = true
   }
   var cnt int = n
   // temporary arrays
   parent := make([]int, n+1)
   size := make([]int, n+1)

   // compute subtree sizes and parents in current alive tree
   var dfs1 func(u, p int)
   dfs1 = func(u, p int) {
       parent[u] = p
       size[u] = 1
       for _, v := range adj[u] {
           if !alive[v] || v == p {
               continue
           }
           dfs1(v, u)
           size[u] += size[v]
       }
   }

   // collect component from start, blocking edge x-y
   var dfs2 func(u, p, x, y int, comp *[]int)
   dfs2 = func(u, p, x, y int, comp *[]int) {
       *comp = append(*comp, u)
       for _, v := range adj[u] {
           if !alive[v] || v == p {
               continue
           }
           // block crossing removed edge
           if (u == x && v == y) || (u == y && v == x) {
               continue
           }
           dfs2(v, u, x, y, comp)
       }
   }

   for cnt > 1 {
       // pick arbitrary root
       var root int
       for i := 1; i <= n; i++ {
           if alive[i] {
               root = i
               break
           }
       }
       dfs1(root, 0)
       // find best edge to query: minimizes max(component sizes)
       bestU, bestV, bestBal := 0, 0, n
       // total alive is size[root]
       total := size[root]
       for u := 1; u <= n; u++ {
           v := parent[u]
           if v == 0 {
               continue
           }
           // u-v edge
           s1 := size[u]
           if s1 == 0 {
               continue
           }
           s2 := total - s1
           bal := s1
           if s2 > bal {
               bal = s2
           }
           if bal < bestBal {
               bestBal = bal
               bestU, bestV = u, v
           }
       }
       // query this edge
       fmt.Fprintln(out, "?", bestU, bestV)
       flush()
       var ans int
       fmt.Fscan(in, &ans)
       // build new alive set: component containing ans, without edge bestU-bestV
       comp := make([]int, 0, total)
       dfs2(ans, 0, bestU, bestV, &comp)
       // reset alive
       for i := 1; i <= n; i++ {
           alive[i] = false
       }
       cnt = 0
       for _, u := range comp {
           alive[u] = true
           cnt++
       }
   }
   // only one remains
   var res int
   for i := 1; i <= n; i++ {
       if alive[i] {
           res = i
           break
       }
   }
   fmt.Fprintln(out, "!", res)
}
