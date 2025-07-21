package main

import (
   "bufio"
   "fmt"
   "os"
)

// Hopcroft-Karp for bipartite matching on nodes 1..n
type HopcroftKarp struct {
   n      int
   adj    [][]int
   pairU  []int
   pairV  []int
   dist   []int
   inf    int
}

func NewHK(n int, adj [][]int) *HopcroftKarp {
   hk := &HopcroftKarp{
       n:     n,
       adj:   adj,
       pairU: make([]int, n+1),
       pairV: make([]int, n+1),
       dist:  make([]int, n+1),
       inf:   1e9,
   }
   return hk
}

// BFS builds layers, returns true if there is an augmenting path
func (hk *HopcroftKarp) bfs(center int) bool {
   queue := make([]int, 0, hk.n)
   for u := 1; u <= hk.n; u++ {
       if u == center {
           hk.dist[u] = hk.inf
       } else if hk.pairU[u] == 0 {
           hk.dist[u] = 0
           queue = append(queue, u)
       } else {
           hk.dist[u] = hk.inf
       }
   }
   hk.dist[0] = hk.inf
   for i := 0; i < len(queue); i++ {
       u := queue[i]
       if hk.dist[u] < hk.dist[0] {
           for _, w := range hk.adj[u] {
               if w == center {
                   continue
               }
               v := hk.pairV[w]
               if hk.dist[v] == hk.inf {
                   hk.dist[v] = hk.dist[u] + 1
                   queue = append(queue, v)
               }
           }
       }
   }
   return hk.dist[0] != hk.inf
}

// DFS searches for augmenting path
func (hk *HopcroftKarp) dfs(u, center int) bool {
   for _, w := range hk.adj[u] {
       if w == center {
           continue
       }
       v := hk.pairV[w]
       if hk.dist[v] == hk.dist[u]+1 {
           if v == 0 || hk.dfs(v, center) {
               hk.pairU[u] = w
               hk.pairV[w] = u
               return true
           }
       }
   }
   hk.dist[u] = hk.inf
   return false
}

// MaxMatching returns max cardinality matching excluding center
func (hk *HopcroftKarp) MaxMatching(center int) int {
   // reset pairs
   for i := 1; i <= hk.n; i++ {
       hk.pairU[i] = 0
       hk.pairV[i] = 0
   }
   matching := 0
   for hk.bfs(center) {
       for u := 1; u <= hk.n; u++ {
           if u != center && hk.pairU[u] == 0 {
               if hk.dfs(u, center) {
                   matching++
               }
           }
       }
   }
   return matching
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   hasEdge := make([][]bool, n+1)
   for i := 1; i <= n; i++ {
       hasEdge[i] = make([]bool, n+1)
   }
   adj := make([][]int, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       if !hasEdge[u][v] {
           hasEdge[u][v] = true
           adj[u] = append(adj[u], v)
       }
   }
   hk := NewHK(n, adj)
   best := m + (3*n - 2) // initial assume k=0: cost = m + |A'| - 2*0
   for center := 1; center <= n; center++ {
       // count center edges present
       cv := 0
       for u := 1; u <= n; u++ {
           if u == center {
               continue
           }
           if hasEdge[u][center] {
               cv++
           }
           if hasEdge[center][u] {
               cv++
           }
       }
       if hasEdge[center][center] {
           cv++
       }
       pv := hk.MaxMatching(center)
       k := cv + pv
       cost := m + (3*n - 2) - 2*k
       if cost < best {
           best = cost
       }
   }
   fmt.Println(best)
}
