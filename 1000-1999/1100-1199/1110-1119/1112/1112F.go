package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   c := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &c[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // build rooted tree at 1
   parent := make([]int, n+1)
   children := make([][]int, n+1)
   order := make([]int, 0, n)
   // stack for DFS
   st := []int{1}
   parent[1] = -1
   for len(st) > 0 {
       v := st[len(st)-1]
       st = st[:len(st)-1]
       order = append(order, v)
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           parent[u] = v
           children[v] = append(children[v], u)
           st = append(st, u)
       }
   }
   const INF = int64(9e18)
   dp0 := make([]int64, n+1)
   dp1 := make([]int64, n+1)
   // C1 and C2
   C1 := make([]int64, n+1)
   C2 := make([]int64, n+1)
   bestDelta := make([]int64, n+1)
   // post-order DP
   for i := n - 1; i >= 0; i-- {
       v := order[i]
       if len(children[v]) == 0 {
           dp0[v] = 0
           dp1[v] = c[v]
           C1[v] = dp1[v]
           C2[v] = INF
           bestDelta[v] = INF
       } else {
           sum0 := int64(0)
           sum1 := int64(0)
           bd := INF
           for _, u := range children[v] {
               sum0 += dp0[u]
               sum1 += dp1[u]
               d := dp0[u] - dp1[u]
               if d < bd {
                   bd = d
               }
           }
           if len(children[v]) <= 1 {
               dp0[v] = sum0
           } else {
               dp0[v] = INF
           }
           C2[v] = sum1
           // C1: purchase v, allow one child zero
           delta := bd
           if delta > 0 {
               delta = 0
           }
           C1[v] = c[v] + sum1 + delta
           // dp1: best of C2 and C1
           if C2[v] < C1[v] {
               dp1[v] = C2[v]
           } else {
               dp1[v] = C1[v]
           }
           bestDelta[v] = bd
       }
   }
   // trace for possible purchases
   avail0 := make([]bool, n+1)
   avail1 := make([]bool, n+1)
   buy := make([]bool, n+1)
   avail1[1] = true
   // BFS from root
   q := []int{1}
   for idx := 0; idx < len(q); idx++ {
       v := q[idx]
       // if in state dp0
       if avail0[v] && dp0[v] < INF {
           for _, u := range children[v] {
               if !avail0[u] {
                   avail0[u] = true
                   q = append(q, u)
               }
           }
       }
       // if in state dp1
       if avail1[v] && dp1[v] < INF {
           // case C2: no purchase v
           if dp1[v] == C2[v] {
               for _, u := range children[v] {
                   if !avail1[u] {
                       avail1[u] = true
                       q = append(q, u)
                   }
               }
           }
           // case C1: purchase v
           if dp1[v] == C1[v] {
               buy[v] = true
               bd := bestDelta[v]
               if len(children[v]) == 0 {
                   // leaf, no children
               } else if bd >= 0 {
                   for _, u := range children[v] {
                       if !avail1[u] {
                           avail1[u] = true
                           q = append(q, u)
                       }
                   }
               } else {
                   // one child can be zero
                   for _, u := range children[v] {
                       d := dp0[u] - dp1[u]
                       if d == bd {
                           if !avail0[u] {
                               avail0[u] = true
                               q = append(q, u)
                           }
                       } else {
                           if !avail1[u] {
                               avail1[u] = true
                               q = append(q, u)
                           }
                       }
                   }
               }
           }
       }
   }
   // collect
   res := make([]int, 0)
   for v := 1; v <= n; v++ {
       if buy[v] {
           res = append(res, v)
       }
   }
   fmt.Fprintln(out, dp1[1], len(res))
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
