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

   var n, m, n1, n2, n3 int
   fmt.Fscan(reader, &n, &m)
   fmt.Fscan(reader, &n1, &n2, &n3)
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   color := make([]int, n+1)
   for i := 1; i <= n; i++ {
       color[i] = -1
   }
   type compSides struct { sides [2][]int }
   var comps []compSides
   // BFS to find bipartite components
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if color[i] != -1 {
           continue
       }
       // new component
       sides := compSides{}
       color[i] = 0
       queue = queue[:0]
       queue = append(queue, i)
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           sides.sides[color[u]] = append(sides.sides[color[u]], u)
           for _, v := range adj[u] {
               if color[v] == -1 {
                   color[v] = 1 - color[u]
                   queue = append(queue, v)
               } else if color[v] == color[u] {
                   fmt.Fprint(writer, "NO")
                   return
               }
           }
       }
       comps = append(comps, sides)
   }
   // DP for selecting group 2 sizes
   compCount := len(comps)
   dp := make([][]bool, compCount+1)
   dp[0] = make([]bool, n2+1)
   dp[0][0] = true
   for i := 1; i <= compCount; i++ {
       dp[i] = make([]bool, n2+1)
       a := len(comps[i-1].sides[0])
       b := len(comps[i-1].sides[1])
       for k := 0; k <= n2; k++ {
           if !dp[i-1][k] {
               continue
           }
           if k+a <= n2 {
               dp[i][k+a] = true
           }
           if k+b <= n2 {
               dp[i][k+b] = true
           }
       }
   }
   if !dp[compCount][n2] {
       fmt.Fprint(writer, "NO")
       return
   }
   // reconstruct choices
   choice := make([]int, compCount)
   k := n2
   for i := compCount; i >= 1; i-- {
       a := len(comps[i-1].sides[0])
       b := len(comps[i-1].sides[1])
       if k >= a && dp[i-1][k-a] {
           choice[i-1] = 0
           k -= a
       } else {
           choice[i-1] = 1
           k -= b
       }
   }
   // assign groups
   ans := make([]int, n+1)
   // assign group 2
   for i, comp := range comps {
       side := choice[i]
       for _, u := range comp.sides[side] {
           ans[u] = 2
       }
   }
   // assign group 1 and 3
   rem1, rem3 := n1, n3
   for u := 1; u <= n; u++ {
       if ans[u] == 0 {
           if rem1 > 0 {
               ans[u] = 1
               rem1--
           } else {
               ans[u] = 3
               rem3--
           }
       }
   }
   // output
   fmt.Fprintln(writer, "YES")
   for u := 1; u <= n; u++ {
       fmt.Fprint(writer, ans[u])
   }
}
