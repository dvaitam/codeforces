package main

import (
   "bufio"
   "fmt"
   "os"
)

// HopcroftKarp for bipartite matching
type HopcroftKarp struct {
   n, m int
   graph [][]int
   pairU, pairV []int
   dist []int
}

func NewHopcroftKarp(n, m int, graph [][]int) *HopcroftKarp {
   return &HopcroftKarp{
       n: n, m: m, graph: graph,
       pairU: make([]int, n+1), pairV: make([]int, m+1), dist: make([]int, n+1),
   }
}

func (hk *HopcroftKarp) bfs() bool {
   const inf = 1<<30
   queue := make([]int, 0, hk.n)
   for u := 1; u <= hk.n; u++ {
       if hk.pairU[u] == 0 {
           hk.dist[u] = 0
           queue = append(queue, u)
       } else {
           hk.dist[u] = inf
       }
   }
   found := false
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range hk.graph[u] {
           if hk.pairV[v] == 0 {
               found = true
           } else if hk.dist[hk.pairV[v]] == inf {
               hk.dist[hk.pairV[v]] = hk.dist[u] + 1
               queue = append(queue, hk.pairV[v])
           }
       }
   }
   return found
}

func (hk *HopcroftKarp) dfs(u int) bool {
   const inf = 1<<30
   for _, v := range hk.graph[u] {
       if hk.pairV[v] == 0 || (hk.dist[hk.pairV[v]] == hk.dist[u]+1 && hk.dfs(hk.pairV[v])) {
           hk.pairU[u] = v
           hk.pairV[v] = u
           return true
       }
   }
   hk.dist[u] = inf
   return false
}

// MaxMatching returns size of maximum matching
func (hk *HopcroftKarp) MaxMatching() int {
   result := 0
   for hk.bfs() {
       for u := 1; u <= hk.n; u++ {
           if hk.pairU[u] == 0 && hk.dfs(u) {
               result++
           }
       }
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   hID := make([][]int, n)
   for i := range hID {
       hID[i] = make([]int, m)
   }
   vID := make([][]int, n)
   for i := range vID {
       vID[i] = make([]int, m)
   }
   hCount := 0
   for i := 0; i < n; i++ {
       j := 0
       for j < m {
           if grid[i][j] == '#' {
               hCount++
               k := j
               for k < m && grid[i][k] == '#' {
                   hID[i][k] = hCount
                   k++
               }
               j = k
           } else {
               j++
           }
       }
   }
   vCount := 0
   for j := 0; j < m; j++ {
       i := 0
       for i < n {
           if grid[i][j] == '#' {
               vCount++
               k := i
               for k < n && grid[k][j] == '#' {
                   vID[k][j] = vCount
                   k++
               }
               i = k
           } else {
               i++
           }
       }
   }
   graph := make([][]int, hCount+1)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '#' {
               u := hID[i][j]
               v := vID[i][j]
               graph[u] = append(graph[u], v)
           }
       }
   }
   hk := NewHopcroftKarp(hCount, vCount, graph)
   ans := hk.MaxMatching()
   fmt.Println(ans)
}
