package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
)

type edge struct {
   to   int
   cost int64
}

// Item for priority queue
type item struct {
   v    int
   dist int64
}

type pq []item

func (h pq) Len() int { return len(h) }
func (h pq) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h pq) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *pq) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *pq) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func dijkstra(start int, adj [][]edge) []int64 {
   n := len(adj)
   const inf = int64(4e18)
   dist := make([]int64, n)
   for i := range dist {
       dist[i] = inf
   }
   dist[start] = 0
   h := &pq{{v: start, dist: 0}}
   heap.Init(h)
   for h.Len() > 0 {
       it := heap.Pop(h).(item)
       v := it.v
       d := it.dist
       if d != dist[v] {
           continue
       }
       for _, e := range adj[v] {
           nd := d + e.cost
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(h, item{v: e.to, dist: nd})
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err == io.EOF {
       return
   }
   var x, y int
   fmt.Fscan(in, &x, &y)
   x--
   y--
   orig := make([][]edge, n)
   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(in, &u, &v, &w)
       u--
       v--
       orig[u] = append(orig[u], edge{to: v, cost: w})
       orig[v] = append(orig[v], edge{to: u, cost: w})
   }
   t := make([]int64, n)
   c := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &t[i], &c[i])
   }
   // build taxi graph
   taxi := make([][]edge, n)
   const inf = int64(4e18)
   for i := 0; i < n; i++ {
       dist := dijkstra(i, orig)
       for j := 0; j < n; j++ {
           if i != j && dist[j] <= t[i] {
               taxi[i] = append(taxi[i], edge{to: j, cost: c[i]})
           }
       }
   }
   // dijkstra on taxi graph from x to y
   dist2 := make([]int64, n)
   for i := range dist2 {
       dist2[i] = inf
   }
   dist2[x] = 0
   h2 := &pq{{v: x, dist: 0}}
   heap.Init(h2)
   for h2.Len() > 0 {
       it := heap.Pop(h2).(item)
       v := it.v
       d := it.dist
       if d != dist2[v] {
           continue
       }
       if v == y {
           break
       }
       for _, e := range taxi[v] {
           nd := d + e.cost
           if nd < dist2[e.to] {
               dist2[e.to] = nd
               heap.Push(h2, item{v: e.to, dist: nd})
           }
       }
   }
   if dist2[y] >= inf/2 {
       fmt.Println(-1)
   } else {
       fmt.Println(dist2[y])
   }
}
