package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // edges and degrees
   type edge struct{ x, y int; w int64 }
   edges := make([]edge, 0, m)
   deg := make([]int, n+1)
   var sumW int64
   // adjacency for connectivity
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       var w int64
       fmt.Fscan(in, &x, &y, &w)
       edges = append(edges, edge{x, y, w})
       sumW += w
       if x == y {
           deg[x] += 2
       } else {
           deg[x]++
           deg[y]++
           adj[x] = append(adj[x], y)
           adj[y] = append(adj[y], x)
       }
       if x == y {
           // loop doesn't affect connectivity beyond the vertex itself
           if len(adj[x]) == 0 {
               adj[x] = append(adj[x], y)
           }
       }
   }
   // BFS from 1 to check reachability of vertices with edges
   vis := make([]bool, n+1)
   queue := []int{1}
   vis[1] = true
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range adj[u] {
           if !vis[v] {
               vis[v] = true
               queue = append(queue, v)
           }
       }
   }
   for v := 1; v <= n; v++ {
       if deg[v] > 0 && !vis[v] {
           fmt.Println(-1)
           return
       }
   }
   // build distance matrix
   const INF = int64(4e18)
   dist := make([][]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = make([]int64, n+1)
       for j := 1; j <= n; j++ {
           if i == j {
               dist[i][j] = 0
           } else {
               dist[i][j] = INF
           }
       }
   }
   for _, e := range edges {
       if e.w < dist[e.x][e.y] {
           dist[e.x][e.y] = e.w
           dist[e.y][e.x] = e.w
       }
   }
   // floyd-warshall
   for k := 1; k <= n; k++ {
       for i := 1; i <= n; i++ {
           if dist[i][k] == INF {
               continue
           }
           for j := 1; j <= n; j++ {
               nd := dist[i][k] + dist[k][j]
               if nd < dist[i][j] {
                   dist[i][j] = nd
               }
           }
       }
   }
   // collect odd degree vertices
   odds := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if deg[i]%2 != 0 {
           odds = append(odds, i)
       }
   }
   k := len(odds)
   if k == 0 {
       fmt.Println(sumW)
       return
   }
   // dp over subsets of odds
   full := 1 << k
   dp := make([]int64, full)
   for mask := 1; mask < full; mask++ {
       dp[mask] = INF
   }
   dp[0] = 0
   for mask := 1; mask < full; mask++ {
       // only even count masks
       if bits := bitsCount(mask); bits%2 != 0 {
           continue
       }
       // find first set bit
       var i int
       for i = 0; i < k; i++ {
           if mask&(1<<i) != 0 {
               break
           }
       }
       // pair i with j
       for j := i + 1; j < k; j++ {
           if mask&(1<<j) != 0 {
               m2 := mask ^ (1<<i) ^ (1<<j)
               cost := dp[m2] + dist[odds[i]][odds[j]]
               if cost < dp[mask] {
                   dp[mask] = cost
               }
           }
       }
   }
   res := sumW + dp[full-1]
   fmt.Println(res)
}

// bitsCount returns population count of bits in x
func bitsCount(x int) int {
   cnt := 0
   for x != 0 {
       x &= x - 1
       cnt++
   }
   return cnt
}
