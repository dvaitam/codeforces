package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   to int
   w  int64
}

type Item struct {
   v    int
   dist int64
}

// Priority queue
type PQ []Item

func (h PQ) Len() int            { return len(h) }
func (h PQ) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h PQ) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PQ) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *PQ) Pop() interface{} {
   old := *h
   n := len(old)
   it := old[n-1]
   *h = old[0 : n-1]
   return it
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int64
   var m int
   var sx, sy, fx, fy int64
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &sx, &sy, &fx, &fy)
   type Point struct{ x, y int64; idx int }
   pts := make([]Point, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &pts[i].x, &pts[i].y)
       pts[i].idx = i
   }
   // Direct walk
   ans := abs64(sx-fx) + abs64(sy-fy)
   if m == 0 {
       fmt.Println(ans)
       return
   }
   // Build graph
   g := make([][]Edge, m)
   // sort by x
   sxord := make([]Point, m)
   syord := make([]Point, m)
   copy(sxord, pts)
   copy(syord, pts)
   sort.Slice(sxord, func(i, j int) bool { return sxord[i].x < sxord[j].x })
   sort.Slice(syord, func(i, j int) bool { return syord[i].y < syord[j].y })
   for i := 0; i < m-1; i++ {
       u := sxord[i].idx
       v := sxord[i+1].idx
       w := sxord[i+1].x - sxord[i].x
       g[u] = append(g[u], Edge{v, w})
       g[v] = append(g[v], Edge{u, w})
       u = syord[i].idx
       v = syord[i+1].idx
       w = syord[i+1].y - syord[i].y
       g[u] = append(g[u], Edge{v, w})
       g[v] = append(g[v], Edge{u, w})
   }
   // Dijkstra
   const INF = int64(9e18)
   dist := make([]int64, m)
   for i := range dist {
       dist[i] = INF
   }
   pq := &PQ{}
   heap.Init(pq)
   // initial distances from start to teleporter
   for i, p := range pts {
       d := min64(abs64(sx-p.x), abs64(sy-p.y))
       dist[i] = d
       heap.Push(pq, Item{i, d})
   }
   // run
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u, d := it.v, it.dist
       if d != dist[u] {
           continue
       }
       // update answer via walking from this teleporter to finish
       cand := d + abs64(pts[u].x-fx) + abs64(pts[u].y-fy)
       if cand < ans {
           ans = cand
       }
       for _, e := range g[u] {
           nd := d + e.w
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{e.to, nd})
           }
       }
   }
   fmt.Println(ans)
}

func abs64(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}
