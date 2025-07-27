package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct {
   to     int
   w      int64
   cost   int
}
type info struct {
   leaves int64
   w      int64
   cost   int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       var S int64
       fmt.Fscan(in, &n, &S)
       adj := make([][]edge, n+1)
       for i := 0; i < n-1; i++ {
           var u, v int
           var w int64
           var c int
           fmt.Fscan(in, &u, &v, &w, &c)
           adj[u] = append(adj[u], edge{v, w, c})
           adj[v] = append(adj[v], edge{u, w, c})
       }
       var edges []info
       var total int64
       // dfs to compute leaves for each edge
       var dfs func(u, p int) int64
       dfs = func(u, p int) int64 {
           if len(adj[u]) == 1 && p != 0 {
               return 1
           }
           var sumLeaves int64
           for _, e := range adj[u] {
               if e.to == p {
                   continue
               }
               cnt := dfs(e.to, u)
               total += e.w * cnt
               edges = append(edges, info{cnt, e.w, e.cost})
               sumLeaves += cnt
           }
           return sumLeaves
       }
       // special case: root with no children, but n>=2 so root always has children
       dfs(1, 0)
       // if already ok
       if total <= S {
           fmt.Fprintln(out, 0)
           continue
       }
       // collect deltas
       var d1, d2 []int64
       for _, ei := range edges {
           w := ei.w
           leaves := ei.leaves
           for w > 0 {
               delta := (w - w/2) * leaves
               if ei.cost == 1 {
                   d1 = append(d1, delta)
               } else {
                   d2 = append(d2, delta)
               }
               w /= 2
           }
       }
       sort.Slice(d1, func(i, j int) bool { return d1[i] > d1[j] })
       sort.Slice(d2, func(i, j int) bool { return d2[i] > d2[j] })
       // prefix sums
       p1 := make([]int64, len(d1)+1)
       for i := 0; i < len(d1); i++ {
           p1[i+1] = p1[i] + d1[i]
       }
       p2 := make([]int64, len(d2)+1)
       for i := 0; i < len(d2); i++ {
           p2[i+1] = p2[i] + d2[i]
       }
       need := total - S
       ans := int64(1e18)
       // try using k2 moves from d2
       for k2 := 0; k2 < len(p2); k2++ {
           sum2 := p2[k2]
           if sum2 >= need {
               cost := int64(k2) * 2
               if cost < ans {
                   ans = cost
               }
               break // more k2 only increases cost
           }
           rem := need - sum2
           // find minimal k1 with p1[k1] >= rem
           // search in p1
           idx := sort.Search(len(p1), func(i int) bool { return p1[i] >= rem })
           if idx < len(p1) {
               cost := int64(k2)*2 + int64(idx)
               if cost < ans {
                   ans = cost
               }
           }
       }
       fmt.Fprintln(out, ans)
   }
}
