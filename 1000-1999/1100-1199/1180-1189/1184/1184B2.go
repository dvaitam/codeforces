package main

import (
   "bufio"
   "container/list"
   "fmt"
   "os"
)

const INF = 1 << 60

// Hopcroft-Karp for bipartite matching
type HopcroftKarp struct {
   n, m int
   adj  [][]int
   dist []int
   pairU []int
   pairV []int
}

func NewHopcroftKarp(n, m int) *HopcroftKarp {
   hk := &HopcroftKarp{
       n: n,
       m: m,
       adj: make([][]int, n+1),
       dist: make([]int, n+1),
       pairU: make([]int, n+1),
       pairV: make([]int, m+1),
   }
   return hk
}

func (hk *HopcroftKarp) AddEdge(u, v int) {
   hk.adj[u] = append(hk.adj[u], v)
}

func (hk *HopcroftKarp) bfs() bool {
   Q := list.New()
   for u := 1; u <= hk.n; u++ {
       if hk.pairU[u] == 0 {
           hk.dist[u] = 0
           Q.PushBack(u)
       } else {
           hk.dist[u] = INF
       }
   }
   dist0 := INF
   for Q.Len() > 0 {
       e := Q.Front()
       u := e.Value.(int)
       Q.Remove(e)
       if hk.dist[u] < dist0 {
           for _, v := range hk.adj[u] {
               pu := hk.pairV[v]
               if pu == 0 {
                   dist0 = hk.dist[u] + 1
               } else if hk.dist[pu] == INF {
                   hk.dist[pu] = hk.dist[u] + 1
                   Q.PushBack(pu)
               }
           }
       }
   }
   return dist0 != INF
}

func (hk *HopcroftKarp) dfs(u int, dist0 int) bool {
   for _, v := range hk.adj[u] {
       pu := hk.pairV[v]
       if pu == 0 || (hk.dist[pu] == hk.dist[u]+1 && hk.dfs(pu, dist0)) {
           hk.pairU[u] = v
           hk.pairV[v] = u
           return true
       }
   }
   hk.dist[u] = INF
   return false
}

func (hk *HopcroftKarp) MaxMatching() int {
   matching := 0
   for hk.bfs() {
       for u := 1; u <= hk.n; u++ {
           if hk.pairU[u] == 0 {
               if hk.dfs(u, 0) {
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
   graph := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       graph[u] = append(graph[u], v)
       graph[v] = append(graph[v], u)
   }
   // compute all-pairs shortest paths via BFS from each node
   dist := make([][]int, n)
   for i := 0; i < n; i++ {
       d := make([]int, n)
       for j := 0; j < n; j++ {
           d[j] = -1
       }
       queue := list.New()
       d[i] = 0
       queue.PushBack(i)
       for queue.Len() > 0 {
           e := queue.Front()
           u := e.Value.(int)
           queue.Remove(e)
           for _, v := range graph[u] {
               if d[v] < 0 {
                   d[v] = d[u] + 1
                   queue.PushBack(v)
               }
           }
       }
       dist[i] = d
   }
   var s, b int
   var k, h int64
   fmt.Fscan(reader, &s, &b, &k, &h)
   ships := make([]struct{x int; a, f int64}, s)
   for i := 0; i < s; i++ {
       fmt.Fscan(reader, &ships[i].x, &ships[i].a, &ships[i].f)
       ships[i].x--
   }
   bases := make([]struct{ x int; d int64 }, b)
   for i := 0; i < b; i++ {
       fmt.Fscan(reader, &bases[i].x, &bases[i].d)
       bases[i].x--
   }
   // build bipartite graph
   hk := NewHopcroftKarp(s, b)
   for i := 0; i < s; i++ {
       for j := 0; j < b; j++ {
           if ships[i].a >= bases[j].d {
               dd := dist[ships[i].x][bases[j].x]
               if dd >= 0 && int64(dd) <= ships[i].f {
                   hk.AddEdge(i+1, j+1)
               }
           }
       }
   }
   m0 := hk.MaxMatching()
   // minimal cost
   cost1 := int64(m0) * k
   cost2 := int64(s) * h
   if cost2 < cost1 {
       fmt.Println(cost2)
   } else {
       fmt.Println(cost1)
   }
}
