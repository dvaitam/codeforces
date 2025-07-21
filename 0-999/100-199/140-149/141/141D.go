package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   to   int
   cost int64
   ramp int // 1-based ramp index, 0 if walking
}

type Item struct {
   node int
   dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) {
   *pq = append(*pq, x.(Item))
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   item := old[n-1]
   *pq = old[0 : n-1]
   return item
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   var L int64
   fmt.Fscan(in, &n, &L)
   type Ramp struct{ x, d, t, p int64; idx int }
   ramps := make([]Ramp, 0, n)
   coords := make([]int64, 0, 3*n+2)
   coords = append(coords, 0, L)
   for i := 1; i <= n; i++ {
       var x, d, t, p int64
       fmt.Fscan(in, &x, &d, &t, &p)
       if x-p < 0 || x+d > L {
           // skip unusable beyond negative or beyond finish
           continue
       }
       ramps = append(ramps, Ramp{x, d, t, p, i})
       coords = append(coords, x-p, x, x+d)
   }
   sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
   coords = unique(coords)
   m := len(coords)
   // map coord to index
   idxOf := func(v int64) int {
       i := sort.Search(len(coords), func(i int) bool { return coords[i] >= v })
       return i
   }
   // build graph
   g := make([][]Edge, m)
   // walking edges
   for i := 0; i+1 < m; i++ {
       w := coords[i+1] - coords[i]
       g[i] = append(g[i], Edge{i + 1, w, 0})
       g[i+1] = append(g[i+1], Edge{i, w, 0})
   }
   // ramp edges
   for _, r := range ramps {
       u := idxOf(r.x - r.p)
       v := idxOf(r.x + r.d)
       cost := r.p + r.t
       g[u] = append(g[u], Edge{v, cost, r.idx})
   }
   // dijkstra
   const INF = int64(4e18)
   dist := make([]int64, m)
   prevNode := make([]int, m)
   prevRamp := make([]int, m)
   for i := range dist {
       dist[i] = INF
       prevNode[i] = -1
   }
   src := idxOf(0)
   dst := idxOf(L)
   dist[src] = 0
   pq := &PriorityQueue{{src, 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u := it.node
       if it.dist != dist[u] {
           continue
       }
       if u == dst {
           break
       }
       for _, e := range g[u] {
           nd := it.dist + e.cost
           if nd < dist[e.to] {
               dist[e.to] = nd
               prevNode[e.to] = u
               prevRamp[e.to] = e.ramp
               heap.Push(pq, Item{e.to, nd})
           }
       }
   }
   // reconstruct path
   path := make([]int, 0)
   for u := dst; u != src; u = prevNode[u] {
       if u < 0 {
           break
       }
       if prevRamp[u] != 0 {
           path = append(path, prevRamp[u])
       }
   }
   // reverse
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   // output
   fmt.Fprintln(out, dist[dst])
   fmt.Fprintln(out, len(path))
   if len(path) > 0 {
       for i, v := range path {
           if i > 0 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, v)
       }
       out.WriteByte('\n')
   } else {
       out.WriteByte('\n')
   }
}

func unique(a []int64) []int64 {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
