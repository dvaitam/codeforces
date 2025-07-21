package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m int
   var k, p int64
   if _, err := fmt.Fscan(reader, &m, &k, &p); err != nil {
       return
   }
   adj := make([][]int, m+1)
   for i := 0; i < m-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS to compute depths
   depth := make([]int, m+1)
   q := make([]int, 0, m)
   q = append(q, 1)
   parent := make([]int, m+1)
   parent[1] = -1
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           depth[v] = depth[u] + 1
           q = append(q, v)
       }
   }
   // collect depths excluding root
   n := m - 1
   depths := make([]int, 0, n)
   for i := 2; i <= m; i++ {
       depths = append(depths, depth[i])
   }
   sort.Ints(depths)
   // prefix sums of depths
   prefix := make([]int64, n+1)
   for i, d := range depths {
       prefix[i+1] = prefix[i] + int64(d)
   }
   // limit k
   if k > int64(n) {
       k = int64(n)
   }
   // binary search maximum S
   lo, hi := int64(0), k
   for lo < hi {
       mid := (lo + hi + 1) / 2
       if can(mid, depths, prefix, p) {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   fmt.Fprintln(writer, lo)
}

// can check if we can select S nodes within budget p
func can(S int64, depths []int, prefix []int64, p int64) bool {
   if S == 0 {
       return true
   }
   n := int64(len(depths))
   // S must be <= n
   if S > n {
       return false
   }
   // slide over possible T at depths[i]
   // depths is sorted ascending, 0-based
   minCost := int64(1<<62 - 1)
   // i is 1-based index in prefix, corresponds to depths[i-1]
   for i := S; i <= n; i++ {
       // consider T = depths[i-1]
       // sum of largest S depths among first i depths: depths[i-S]..depths[i-1]
       sumSeg := prefix[i] - prefix[i-S]
       T := int64(depths[i-1])
       cost := S*T - sumSeg
       if cost < minCost {
           minCost = cost
       }
       // early exit if within budget
       if minCost <= p {
           return true
       }
   }
   return minCost <= p
}
