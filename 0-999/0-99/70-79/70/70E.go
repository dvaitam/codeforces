package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, k int
   w    []int
   adj  [][]int
   d    [][]int
   f    [][]int
   b    []int
   c    []int
)

func dfs(u, parent int) {
   // initialize dp for u
   for j := 1; j <= n; j++ {
       f[u][j] = w[d[u][j]] + k
   }
   // process children
   for _, v := range adj[u] {
       if v == parent {
           continue
       }
       dfs(v, u)
       for j := 1; j <= n; j++ {
           // either use same center j, saving k cost at child, or use child's best
           a := f[v][j] - k
           bch := f[v][ b[v] ]
           if a < bch {
               f[u][j] += a
           } else {
               f[u][j] += bch
           }
       }
   }
   // choose best center for u
   best := 1
   for j := 2; j <= n; j++ {
       if f[u][j] < f[u][best] {
           best = j
       }
   }
   b[u] = best
}

func prt(u, parent, z int) {
   c[u] = z
   for _, v := range adj[u] {
       if v == parent {
           continue
       }
       // decide center for child v
       if f[v][z]-k < f[v][ b[v] ] {
           prt(v, u, z)
       } else {
           prt(v, u, b[v])
       }
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &k)
   w = make([]int, n+1)
   for i := 1; i < n; i++ {
       fmt.Fscan(in, &w[i])
   }
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   // compute all-pairs distances by BFS
   const inf = 1<<30
   d = make([][]int, n+1)
   for u := 1; u <= n; u++ {
       dist := make([]int, n+1)
       for i := 1; i <= n; i++ {
           dist[i] = -1
       }
       queue := make([]int, 0, n)
       dist[u] = 0
       queue = append(queue, u)
       for qi := 0; qi < len(queue); qi++ {
           v := queue[qi]
           for _, to := range adj[v] {
               if dist[to] < 0 {
                   dist[to] = dist[v] + 1
                   queue = append(queue, to)
               }
           }
       }
       row := make([]int, n+1)
       for v := 1; v <= n; v++ {
           if dist[v] >= 0 {
               row[v] = dist[v]
           } else {
               row[v] = inf
           }
       }
       d[u] = row
   }
   // init dp tables
   f = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       f[i] = make([]int, n+1)
   }
   b = make([]int, n+1)
   c = make([]int, n+1)
   // dp
   dfs(1, 0)
   rootBest := b[1]
   // output cost and assignments
   fmt.Println(f[1][rootBest])
   prt(1, 0, rootBest)
   for i := 1; i <= n; i++ {
       if i > 1 {
           fmt.Print(" ")
       }
       fmt.Print(c[i])
   }
   fmt.Println()
}
