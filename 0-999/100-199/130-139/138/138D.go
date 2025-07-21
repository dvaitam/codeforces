package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // diagonals: s = i+j, d = i-j+(m-1)
   S := n + m - 1
   T := S
   l := make([]bool, S)
   r := make([]bool, T)
   adj := make([][]int, S)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           c := grid[i][j]
           u := i + j
           v := i - j + (m - 1)
           switch c {
           case 'L':
               l[u] = true
           case 'R':
               r[v] = true
           case 'X':
               adj[u] = append(adj[u], v)
           }
       }
   }
   // Hopcroft-Karp
   const INF = 1 << 30
   pairU := make([]int, S)
   pairV := make([]int, T)
   dist := make([]int, S)
   for i := range pairU {
       pairU[i] = -1
   }
   for i := range pairV {
       pairV[i] = -1
   }
   var bfs func() bool
   var dfs func(u int) bool
   bfs = func() bool {
       queue := make([]int, 0, S)
       for u := 0; u < S; u++ {
           if pairU[u] < 0 {
               dist[u] = 0
               queue = append(queue, u)
           } else {
               dist[u] = INF
           }
       }
       found := false
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if pairV[v] < 0 {
                   found = true
               } else if dist[pairV[v]] == INF {
                   dist[pairV[v]] = dist[u] + 1
                   queue = append(queue, pairV[v])
               }
           }
       }
       return found
   }
   dfs = func(u int) bool {
       for _, v := range adj[u] {
           if pairV[v] < 0 || (dist[pairV[v]] == dist[u]+1 && dfs(pairV[v])) {
               pairU[u] = v
               pairV[v] = u
               return true
           }
       }
       dist[u] = INF
       return false
   }
   matching := 0
   for bfs() {
       for u := 0; u < S; u++ {
           if pairU[u] < 0 && dfs(u) {
               matching++
           }
       }
   }
   countL, countR := 0, 0
   for u := 0; u < S; u++ {
       if l[u] {
           countL++
       }
   }
   for v := 0; v < T; v++ {
       if r[v] {
           countR++
       }
   }
   val := countL + countR - 2*matching
   if val&1 != 0 {
       fmt.Println("WIN")
   } else {
       fmt.Println("LOSE")
   }
}
