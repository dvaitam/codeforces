package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   const INF = 1e18
   // distance matrix
   dist := make([][]int64, n)
   for i := 0; i < n; i++ {
       dist[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           if i == j {
               dist[i][j] = 0
           } else {
               dist[i][j] = INF
           }
       }
   }
   for i := 0; i < m; i++ {
       var u, v int
       var c int64
       fmt.Fscan(reader, &u, &v, &c)
       u--
       v--
       if c < dist[u][v] {
           dist[u][v] = c
           dist[v][u] = c
       }
   }
   // Floyd-Warshall
   for k := 0; k < n; k++ {
       for i := 0; i < n; i++ {
           if dist[i][k] == INF {
               continue
           }
           for j := 0; j < n; j++ {
               if dist[k][j] == INF {
                   continue
               }
               nd := dist[i][k] + dist[k][j]
               if nd < dist[i][j] {
                   dist[i][j] = nd
               }
           }
       }
   }
   var g1, g2, s1, s2 int
   fmt.Fscan(reader, &g1, &g2, &s1, &s2)
   // collect unique distances
   dvals := make([]int64, 0, n*(n-1))
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if i != j {
               dvals = append(dvals, dist[i][j])
           }
       }
   }
   sort.Slice(dvals, func(i, j int) bool { return dvals[i] < dvals[j] })
   dvals = uniqueInt64(dvals)
   D := len(dvals)
   // A[p][i] = count j: dist[i][j] <= dvals[p]
   A := make([][]int, D)
   for p := 0; p < D; p++ {
       A[p] = make([]int, n)
       for i := 0; i < n; i++ {
           cnt := 0
           for j := 0; j < n; j++ {
               if i != j && dist[i][j] <= dvals[p] {
                   cnt++
               }
           }
           A[p][i] = cnt
       }
   }
   ans := int64(0)
   // DP arrays
   maxG := g2
   maxS := s2
   // iterate boundaries
   for p := 0; p < D; p++ {
       for q := p + 1; q < D; q++ {
           // for each runner, counts for categories
           a := make([]int64, n)
           b := make([]int64, n)
           c := make([]int64, n)
           for i := 0; i < n; i++ {
               a[i] = int64(A[p][i])
               b[i] = int64(A[q][i] - A[p][i])
               c[i] = int64((n - 1) - A[q][i])
           }
           // dp[k][l]
           dp := make([][]int64, maxG+1)
           for i := range dp {
               dp[i] = make([]int64, maxS+1)
           }
           dp[0][0] = 1
           for i := 0; i < n; i++ {
               ndp := make([][]int64, maxG+1)
               for ii := range ndp {
                   ndp[ii] = make([]int64, maxS+1)
               }
               for gi := 0; gi <= maxG; gi++ {
                   for si := 0; si <= maxS; si++ {
                       v := dp[gi][si]
                       if v == 0 {
                           continue
                       }
                       // gold
                       if gi < maxG && a[i] > 0 {
                           ndp[gi+1][si] += v * a[i]
                       }
                       // silver
                       if si < maxS && b[i] > 0 {
                           ndp[gi][si+1] += v * b[i]
                       }
                       // bronze
                       if c[i] > 0 {
                           ndp[gi][si] += v * c[i]
                       }
                   }
               }
               dp = ndp
           }
           // sum valid
           for gi := g1; gi <= g2; gi++ {
               if gi > maxG {
                   break
               }
               for si := s1; si <= s2; si++ {
                   if si > maxS {
                       break
                   }
                   ans += dp[gi][si]
               }
           }
       }
   }
   fmt.Println(ans)
}

func uniqueInt64(a []int64) []int64 {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
