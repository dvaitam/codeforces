package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

const INF64 = (1<<63 - 1) / 2

// For simple Dijkstra from source
func dijkstra(n int, adj [][]edge, src int) []int64 {
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF64
   }
   dist[src] = 0
   pq := &hp{}
   heap.Init(pq)
   heap.Push(pq, item{u: src, d: 0})
   for pq.Len() > 0 {
       cur := heap.Pop(pq).(item)
       u, d := cur.u, cur.d
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           nd := d + e.w
           if nd < dist[e.v] {
               dist[e.v] = nd
               heap.Push(pq, item{u: e.v, d: nd})
           }
       }
   }
   return dist
}

// Multi-source Dijkstra: returns dist and label of nearest source
func multiSource(n int, adj [][]edge, sources []int) ([]int64, []int) {
   dist := make([]int64, n+1)
   label := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF64
       label[i] = 0
   }
   pq := &msHP{}
   heap.Init(pq)
   for _, s := range sources {
       dist[s] = 0
       label[s] = s
       heap.Push(pq, msItem{u: s, d: 0})
   }
   for pq.Len() > 0 {
       cur := heap.Pop(pq).(msItem)
       u, d := cur.u, cur.d
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           nd := d + e.w
           if nd < dist[e.v] {
               dist[e.v] = nd
               label[e.v] = label[u]
               heap.Push(pq, msItem{u: e.v, d: nd})
           }
       }
   }
   return dist, label
}

// edge for graph
type edge struct { v int; w int64 }
// original edge
type origEdge struct { u, v int; w int64 }

// item for Dijkstra
type item struct { u int; d int64 }
// min-heap of items
type hp []item
func (h hp) Len() int            { return len(h) }
func (h hp) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h hp) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *hp) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

// msItem for multi-source
type msItem struct { u int; d int64 }
// heap for ms items
type msHP []msItem
func (h msHP) Len() int           { return len(h) }
func (h msHP) Less(i, j int) bool { return h[i].d < h[j].d }
func (h msHP) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *msHP) Push(x interface{}) { *h = append(*h, x.(msItem)) }
func (h *msHP) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

// DSU
type dsu struct { p []int }
func newDSU(n int) *dsu {
   p := make([]int, n+1)
   for i := range p {
       p[i] = i
   }
   return &dsu{p: p}
}
func (d *dsu) find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.find(d.p[x])
   }
   return d.p[x]
}
func (d *dsu) union(x, y int) bool {
   rx, ry := d.find(x), d.find(y)
   if rx == ry {
       return false
   }
   d.p[ry] = rx
   return true
}

// edge for MST
type mstEdge struct { u, v int; w int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   adj := make([][]edge, n+1)
   orig := make([]origEdge, 0, m)
   for i := 0; i < m; i++ {
       var x, y int
       var w int64
       fmt.Fscan(in, &x, &y, &w)
       adj[x] = append(adj[x], edge{v: y, w: w})
       adj[y] = append(adj[y], edge{v: x, w: w})
       orig = append(orig, origEdge{u: x, v: y, w: w})
   }
   var k int
   fmt.Fscan(in, &k)
   ps := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(in, &ps[i])
   }
   // Dijkstra from city 1
   dist1 := dijkstra(n, adj, 1)
   // multi-source from portals
   dist, label := multiSource(n, adj, ps)
   // build edges between portals
   edges := make([]mstEdge, 0, m)
   for _, e := range orig {
       lu, lv := label[e.u], label[e.v]
       if lu != lv {
           w2 := dist[e.u] + e.w + dist[e.v]
           edges = append(edges, mstEdge{u: lu, v: lv, w: w2})
       }
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].w < edges[j].w
   })
   dsu := newDSU(n)
   var tot int64
   cnt := 0
   for _, e := range edges {
       if dsu.union(e.u, e.v) {
           tot += e.w
           cnt++
           if cnt == k-1 {
               break
           }
       }
   }
   // add cost from city1 to first portal
   minD1 := INF64
   for _, p := range ps {
       if dist1[p] < minD1 {
           minD1 = dist1[p]
       }
   }
   if minD1 == INF64 {
       minD1 = 0
   }
   ans := tot + minD1
   fmt.Println(ans)
}
