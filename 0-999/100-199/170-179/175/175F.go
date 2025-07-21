package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Edge represents an edge in the graph
type Edge struct {
   to int
}

// State for priority queue
type State struct {
   bad, dist int
   path      []int
   node      int
}

// Priority queue
type PQ []*State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
   a, b := pq[i], pq[j]
   if a.bad != b.bad {
       return a.bad < b.bad
   }
   if a.dist != b.dist {
       return a.dist < b.dist
   }
   // lex compare paths
   pa, pb := a.path, b.path
   na, nb := len(pa), len(pb)
   lim := na
   if nb < lim {
       lim = nb
   }
   for i := 0; i < lim; i++ {
       if pa[i] != pb[i] {
           return pa[i] < pb[i]
       }
   }
   return na < nb
}
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*State)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[:n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   good := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &good[i])
   }
   adj := make([][]Edge, n+1)
   // read shortcuts and main cycle
   // add main cycle edges
   for i := 0; i < m; i++ {
       u := good[i]
       v := good[(i+1)%m]
       adj[u] = append(adj[u], Edge{v})
       adj[v] = append(adj[v], Edge{u})
   }
   // evil shortcuts
   for i := 0; i < m; i++ {
       var k int
       fmt.Fscan(reader, &k)
       pts := make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &pts[j])
       }
       for j := 0; j+1 < k; j++ {
           u, v := pts[j], pts[j+1]
           adj[u] = append(adj[u], Edge{v})
           adj[v] = append(adj[v], Edge{u})
       }
   }
   // sort adjacency for lex
   for i := 1; i <= n; i++ {
       sort.Slice(adj[i], func(a, b int) bool {
           return adj[i][a].to < adj[i][b].to
       })
   }
   // edge weights: map from u*n+v to bad count
   weights := make(map[int]int)
   q := 0
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       var op byte
       var s, t int
       fmt.Fscan(reader, &op, &s, &t)
       if op == '+' {
           key1 := s*(n+1) + t
           key2 := t*(n+1) + s
           weights[key1]++
           weights[key2]++
       } else if op == '?' {
           // Dijkstra with tuple (bad, dist, lex path)
           dist := make([]int, n+1)
           badc := make([]int, n+1)
           vis := make([]bool, n+1)
           paths := make([][]int, n+1)
           const INF = 1<<60
           for j := 1; j <= n; j++ {
               dist[j] = 1e9
               badc[j] = 1e9
           }
           pq := &PQ{}
           heap.Init(pq)
           badc[s], dist[s] = 0, 0
           paths[s] = []int{s}
           heap.Push(pq, &State{0, 0, []int{s}, s})
           var ansState *State
           for pq.Len() > 0 {
               cur := heap.Pop(pq).(*State)
               u := cur.node
               if vis[u] {
                   continue
               }
               vis[u] = true
               if u == t {
                   ansState = cur
                   break
               }
               for _, e := range adj[u] {
                   v := e.to
                   if vis[v] {
                       continue
                   }
                   key := u*(n+1) + v
                   w := weights[key]
                   nb := cur.bad + w
                   nd := cur.dist + 1
                   newPath := append([]int{}, cur.path...)
                   newPath = append(newPath, v)
                   better := false
                   if nb < badc[v] || (nb == badc[v] && nd < dist[v]) {
                       better = true
                   } else if nb == badc[v] && nd == dist[v] {
                       // lex compare
                       pOld := paths[v]
                       for x := 0; x < len(newPath) && x < len(pOld); x++ {
                           if newPath[x] != pOld[x] {
                               if newPath[x] < pOld[x] {
                                   better = true
                               }
                               break
                           }
                       }
                       if !better && len(newPath) < len(pOld) {
                           better = true
                       }
                   }
                   if better {
                       badc[v], dist[v] = nb, nd
                       paths[v] = newPath
                       heap.Push(pq, &State{nb, nd, newPath, v})
                   }
               }
           }
           if ansState == nil {
               fmt.Fprintln(writer, -1)
           } else {
               fmt.Fprintln(writer, ansState.bad)
               // clear weights on path
               for j := 1; j < len(ansState.path); j++ {
                   u := ansState.path[j-1]
                   v := ansState.path[j]
                   key1 := u*(n+1) + v
                   key2 := v*(n+1) + u
                   delete(weights, key1)
                   delete(weights, key2)
               }
           }
       }
   }
}
