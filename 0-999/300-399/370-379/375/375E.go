package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n int
   x int64
   color []int
   adj [][]edge
)

type edge struct { to int; w int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &x)
   color = make([]int, n)
   var bcnt int
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &color[i])
       if color[i] == 1 {
           bcnt++
       }
   }
   adj = make([][]edge, n)
   for i := 0; i < n-1; i++ {
       var u, v int; var w int64
       fmt.Fscan(in, &u, &v, &w)
       u--; v--
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   if bcnt == 0 {
       // no black vertices, cannot cover
       fmt.Println(-1)
       return
   }
   // dp returns map: k -> map[d] -> max intersections in subtree
   dpRoot := dfs(0, -1)
   // find best for k = bcnt
   best := -1
   if m, ok := dpRoot[bcnt]; ok {
       for _, inter := range m {
           if inter > best {
               best = inter
           }
       }
   }
   if best < 0 {
       fmt.Println(-1)
   } else {
       // operations = bcnt - best
       fmt.Println(bcnt - best)
   }
}

// dfs returns dp map for subtree u: dp[k][d] = max intersections
func dfs(u, p int) map[int]map[int64]int {
   // dp0: u has no center; dp1: u has a center
   dp0 := make(map[int]map[int64]int)
   dp1 := make(map[int]map[int64]int)
   // init
   dp0[0] = map[int64]int{0: 0}
   initInter := 0
   if color[u] == 1 {
       initInter = 1
   }
   // center at u covers itself: no uncovered nodes => d = -1
   dp1[1] = map[int64]int{-1: initInter}
   // children
   for _, e := range adj[u] {
       v, w := e.to, e.w
       if v == p {
           continue
       }
       dpv := dfs(v, u)
       // merge to new dp0 and dp1
       new0 := make(map[int]map[int64]int)
       new1 := make(map[int]map[int64]int)
       // merge for dp0
       for k1, m1 := range dp0 {
           for d1, inter1 := range m1 {
               for k2, m2 := range dpv {
                   for d2, inter2 := range m2 {
                       k := k1 + k2
                       // no center at u: propagate uncovered distances
                       var d_eff int64
                       if d2 < 0 {
                           d_eff = d1
                       } else {
                           d2w := d2 + w
                           if d2w < d1 {
                               d_eff = d1
                           } else {
                               d_eff = d2w
                           }
                       }
                       if d_eff <= x {
                           inter := inter1 + inter2
                           mm, ok := new0[k]
                           if !ok {
                               mm = make(map[int64]int)
                               new0[k] = mm
                           }
                           if prev, ok2 := mm[d_eff]; !ok2 || inter > prev {
                               mm[d_eff] = inter
                           }
                       }
                   }
               }
           }
       }
       // merge for dp1 (center at u)
       for k1, m1 := range dp1 {
           for d1, inter1 := range m1 {
               for k2, m2 := range dpv {
                   for d2, inter2 := range m2 {
                       k := k1 + k2
                       // has center at u: child uncovered nodes may be covered by u
                       var d_eff int64
                       if d2 < 0 {
                           // child fully covered
                           d_eff = d1
                       } else {
                           d2w := d2 + w
                           if d2w <= x {
                               // covered by u
                               d_eff = d1
                           } else {
                               // remains uncovered
                               if d2w < d1 {
                                   d_eff = d1
                               } else {
                                   d_eff = d2w
                               }
                           }
                       }
                       if d_eff <= x {
                           inter := inter1 + inter2
                           mm, ok := new1[k]
                           if !ok {
                               mm = make(map[int64]int)
                               new1[k] = mm
                           }
                           if prev, ok2 := mm[d_eff]; !ok2 || inter > prev {
                               mm[d_eff] = inter
                           }
                       }
                   }
               }
           }
       }
       // prune new0 and new1
       dp0 = prune(new0)
       dp1 = prune(new1)
   }
   // combine dp0 and dp1
   dp := make(map[int]map[int64]int)
   for k, m := range dp0 {
       mm := dp[k]
       if mm == nil {
           mm = make(map[int64]int)
       }
       for d, inter := range m {
           if prev, ok := mm[d]; !ok || inter > prev {
               mm[d] = inter
           }
       }
       dp[k] = mm
   }
   for k, m := range dp1 {
       mm := dp[k]
       if mm == nil {
           mm = make(map[int64]int)
       }
       for d, inter := range m {
           if prev, ok := mm[d]; !ok || inter > prev {
               mm[d] = inter
           }
       }
       dp[k] = mm
   }
   // final prune
   return prune(dp)
}

// prune keeps for each k only Pareto-optimal (d, inter)
func prune(dp map[int]map[int64]int) map[int]map[int64]int {
   res := make(map[int]map[int64]int)
   type pair struct{ d, inter int64 }
   for k, m := range dp {
       ps := make([]pair, 0, len(m))
       for d, inter := range m {
           ps = append(ps, pair{d, int64(inter)})
       }
       // sort by d asc, inter desc
       for i := 1; i < len(ps); i++ {
           for j := i; j > 0 && (ps[j].d < ps[j-1].d || (ps[j].d == ps[j-1].d && ps[j].inter > ps[j-1].inter)); j-- {
               ps[j], ps[j-1] = ps[j-1], ps[j]
           }
       }
       mm := make(map[int64]int)
       var maxInter int64 = -1
       for _, pr := range ps {
           if pr.inter > maxInter {
               mm[pr.d] = int(pr.inter)
               maxInter = pr.inter
           }
       }
       res[k] = mm
   }
   return res
}
