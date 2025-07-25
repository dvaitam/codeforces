package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

type Edge struct {
   to int
   w  int
}

// state for Dijkstra: cost to reach node u with k b-edges used
type State struct {
   cost int
   u    int
   k    int
}

// priority queue of State
type PQ []State

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].cost < pq[j].cost }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[0 : n-1]
   return x
}

// union-find
type UF struct { p []int }
func NewUF(n int) *UF {
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &UF{p}
}
func (uf *UF) Find(x int) int {
   if uf.p[x] != x {
       uf.p[x] = uf.Find(uf.p[x])
   }
   return uf.p[x]
}
func (uf *UF) Union(x, y int) bool {
   rx, ry := uf.Find(x), uf.Find(y)
   if rx == ry {
       return false
   }
   uf.p[ry] = rx
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var a, b int
   fmt.Fscan(in, &n, &m, &a, &b)
   type E struct { u, v, w int }
   edges := make([]E, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
   }
   // compute countA: number of a-edges in MST
   uf := NewUF(n)
   countA := 0
   // first take all a-edges
   for _, e := range edges {
       if e.w == a {
           if uf.Union(e.u, e.v) {
               countA++
           }
       }
   }
   // then b-edges to fill MST, but we only need countA
   needB := (n - 1) - countA
   // build adjacency
   adj := make([][]Edge, n+1)
   for _, e := range edges {
       adj[e.u] = append(adj[e.u], Edge{e.v, e.w})
       adj[e.v] = append(adj[e.v], Edge{e.u, e.w})
   }
   const INF = 1<<60
   // dp[u][k] = minimal cost to reach u with k b-edges
   dp := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       dp[i] = make([]int, needB+1)
       for k := 0; k <= needB; k++ {
           dp[i][k] = INF
       }
   }
   dp[1][0] = 0
   pq := &PQ{}
   heap.Init(pq)
   heap.Push(pq, State{0, 1, 0})
   for pq.Len() > 0 {
       st := heap.Pop(pq).(State)
       if st.cost != dp[st.u][st.k] {
           continue
       }
       u := st.u
       ck := st.k
       for _, e := range adj[u] {
           if e.w == a {
               nc := st.cost + a
               if nc < dp[e.to][ck] {
                   dp[e.to][ck] = nc
                   heap.Push(pq, State{nc, e.to, ck})
               }
           } else {
               if ck < needB {
                   nc := st.cost + b
                   if nc < dp[e.to][ck+1] {
                       dp[e.to][ck+1] = nc
                       heap.Push(pq, State{nc, e.to, ck + 1})
                   }
               }
           }
       }
   }
   // output answers
   for i := 1; i <= n; i++ {
       best := INF
       for k := 0; k <= needB; k++ {
           if dp[i][k] < best {
               best = dp[i][k]
           }
       }
       if i > 1 {
           fmt.Fprint(out, ' ')
       }
       fmt.Fprint(out, best)
   }
   fmt.Fprintln(out)
}
