package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

const INF64 = int64(9e18)

// Dijkstra priority queue
type Item struct {
   node int
   dist int64
}
type PQ []Item
func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

// min-heap for int64
type MinHeap struct { data []int64 }
func (h *MinHeap) Len() int { return len(h.data) }
func (h *MinHeap) Less(i, j int) bool { return h.data[i] < h.data[j] }
func (h *MinHeap) Swap(i, j int) { h.data[i], h.data[j] = h.data[j], h.data[i] }
func (h *MinHeap) Push(x interface{}) { h.data = append(h.data, x.(int64)) }
func (h *MinHeap) Pop() interface{} {
   old := h.data
   n := len(old)
   v := old[n-1]
   h.data = old[:n-1]
   return v
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, q int
   fmt.Fscan(in, &n, &m, &q)
   U := make([]int, m+1)
   V := make([]int, m+1)
   W := make([]int64, m+1)
   adj := make([][]struct{to int; w int64; idx int}, n+1)
   for i := 1; i <= m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(in, &u, &v, &w)
       U[i], V[i], W[i] = u, v, w
       adj[u] = append(adj[u], struct{to int; w int64; idx int}{v, w, i})
       adj[v] = append(adj[v], struct{to int; w int64; idx int}{u, w, i})
   }
   // Dijkstra from 1
   distS := make([]int64, n+1)
   for i := range distS { distS[i] = INF64 }
   distS[1] = 0
   parent := make([]int, n+1)
   parentEdge := make([]int, n+1)
   pq := &PQ{{1, 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u, d := it.node, it.dist
       if d != distS[u] { continue }
       for _, e := range adj[u] {
           v, w, idx := e.to, e.w, e.idx
           nd := d + w
           if nd < distS[v] {
               distS[v] = nd
               parent[v] = u
               parentEdge[v] = idx
               heap.Push(pq, Item{v, nd})
           }
       }
   }
   // Dijkstra from n
   distT := make([]int64, n+1)
   for i := range distT { distT[i] = INF64 }
   distT[n] = 0
   pq = &PQ{{n, 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u, d := it.node, it.dist
       if d != distT[u] { continue }
       for _, e := range adj[u] {
           v, w := e.to, e.w
           nd := d + w
           if nd < distT[v] {
               distT[v] = nd
               heap.Push(pq, Item{v, nd})
           }
       }
   }
   // retrieve shortest path P
   var P []int
   var pathEdges []int
   cur := n
   for cur != 1 {
       P = append(P, cur)
       pathEdges = append(pathEdges, parentEdge[cur])
       cur = parent[cur]
   }
   P = append(P, 1)
   // reverse
   for i, j := 0, len(P)-1; i < j; i, j = i+1, j-1 { P[i], P[j] = P[j], P[i] }
   for i, j := 0, len(pathEdges)-1; i < j; i, j = i+1, j-1 { pathEdges[i], pathEdges[j] = pathEdges[j], pathEdges[i] }
   k := len(P)
   idx := make([]int, n+1)
   for i := range idx { idx[i] = -1 }
   for i, u := range P { idx[u] = i }
   // map edge index to path position
   edgePos := make([]int, m+1)
   for i := range edgePos { edgePos[i] = -1 }
   for i, eidx := range pathEdges {
       edgePos[eidx] = i
   }
   // build DAG_S predecessors
   predS := make([][]int, n+1)
   for i := 1; i <= m; i++ {
       u, v, w := U[i], V[i], W[i]
       if distS[u]+w == distS[v] {
           predS[v] = append(predS[v], u)
       }
       if distS[v]+w == distS[u] {
           predS[u] = append(predS[u], v)
       }
   }
   // build DAG_T successors
   succT := make([][]int, n+1)
   for i := 1; i <= m; i++ {
       u, v, w := U[i], V[i], W[i]
       // edges towards n
       if distT[u] == w+distT[v] {
           succT[u] = append(succT[u], v)
       }
       if distT[v] == w+distT[u] {
           succT[v] = append(succT[v], u)
       }
   }
   // posFirst
   ordS := make([]int, n)
   for i := 1; i <= n; i++ { ordS[i-1] = i }
   sort.Slice(ordS, func(i, j int) bool { return distS[ordS[i]] < distS[ordS[j]] })
   posFirst := make([]int, n+1)
   for _, u := range ordS {
       best := -1
       if idx[u] >= 0 {
           best = idx[u]
       }
       for _, p := range predS[u] {
           if posFirst[p] > best {
               best = posFirst[p]
           }
       }
       posFirst[u] = best
   }
   // posLast
   ordT := make([]int, n)
   for i := 1; i <= n; i++ { ordT[i-1] = i }
   // process nodes by increasing distT so successors (closer to n) are first
   sort.Slice(ordT, func(i, j int) bool { return distT[ordT[i]] < distT[ordT[j]] })
   posLast := make([]int, n+1)
   for _, u := range ordT {
       best := k
       if idx[u] >= 0 {
           best = idx[u]
       }
       for _, v := range succT[u] {
           if posLast[v] < best {
               best = posLast[v]
           }
       }
       posLast[u] = best
   }
   // events for sweep
   ins := make([][]int64, k)
   del := make([][]int64, k)
   for i := 1; i <= m; i++ {
       if edgePos[i] != -1 {
           continue
       }
       u, v, w := U[i], V[i], W[i]
       // u->v
       l, r := posFirst[u], posLast[v]
       if l < r {
           cost := distS[u] + w + distT[v]
           ins[l] = append(ins[l], cost)
           del[r] = append(del[r], cost)
       }
       // v->u
       l, r = posFirst[v], posLast[u]
       if l < r {
           cost := distS[v] + w + distT[u]
           ins[l] = append(ins[l], cost)
           del[r] = append(del[r], cost)
       }
   }
   // sweep to compute alt
   alt := make([]int64, k-1)
   mh := &MinHeap{data: make([]int64, 0)}
   heap.Init(mh)
   delMap := make(map[int64]int)
   for i := 0; i < k-1; i++ {
       for _, v := range ins[i] {
           heap.Push(mh, v)
       }
       for _, v := range del[i] {
           delMap[v]++
       }
       // clean
       for mh.Len() > 0 {
           top := mh.data[0]
           if delMap[top] > 0 {
               heap.Pop(mh)
               delMap[top]--
           } else {
               break
           }
       }
       if mh.Len() == 0 {
           alt[i] = INF64
       } else {
           alt[i] = mh.data[0]
       }
   }
   // answer queries
   base := distS[n]
   for qi := 0; qi < q; qi++ {
       var ti int
       var x int64
       fmt.Fscan(in, &ti, &x)
       u, v := U[ti], V[ti]
       orig := W[ti]
       pos := edgePos[ti]
       var ans int64
       if pos == -1 {
           ans = base
           // try using this edge
           c1 := distS[u] + x + distT[v]
           if c1 < ans { ans = c1 }
           c2 := distS[v] + x + distT[u]
           if c2 < ans { ans = c2 }
       } else {
           // on path
           c1 := base + x - orig
           ans = c1
           if alt[pos] < ans {
               ans = alt[pos]
           }
       }
       fmt.Fprintln(out, ans)
   }
}
