package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
   "sort"
)

const inf64 = int64(4e18)

// Pair for max-heap
type Pair struct {
   val float64
   idx int
}
type MaxHeap []Pair
func (h MaxHeap) Len() int            { return len(h) }
func (h MaxHeap) Less(i, j int) bool  { return h[i].val > h[j].val }
func (h MaxHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) { *h = append(*h, x.(Pair)) }
func (h *MaxHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[:n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   edges := make([]struct{u, v int; w int64}, m)
   dist := make([][]int64, n)
   for i := 0; i < n; i++ {
       dist[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           if i == j {
               dist[i][j] = 0
           } else {
               dist[i][j] = inf64
           }
       }
   }
   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(in, &u, &v, &w)
       u--
       v--
       edges[i] = struct{u, v int; w int64}{u, v, w}
       if w < dist[u][v] {
           dist[u][v] = w
           dist[v][u] = w
       }
   }
   // All-pairs shortest paths
   for k := 0; k < n; k++ {
       for i := 0; i < n; i++ {
           if dist[i][k] == inf64 {
               continue
           }
           for j := 0; j < n; j++ {
               d := dist[i][k] + dist[k][j]
               if d < dist[i][j] {
                   dist[i][j] = d
               }
           }
       }
   }
   // best at vertex
   best := math.Inf(1)
   for u := 0; u < n; u++ {
       mx := int64(0)
       for i := 0; i < n; i++ {
           if dist[u][i] > mx {
               mx = dist[u][i]
           }
       }
       d := float64(mx)
       if d < best {
           best = d
       }
   }
   // consider edges
   for _, e := range edges {
       d2 := solveEdge(e.u, e.v, e.w, dist)
       if d2 < best {
           best = d2
       }
   }
   fmt.Printf("%.10f\n", best)
}

// solveEdge returns minimal maximum distance for nodes from a point on edge u-v of weight w
func solveEdge(u, v int, w int64, dist [][]int64) float64 {
   n := len(dist)
   a := make([]float64, n)
   b := make([]float64, n)
   bw := float64(w)
   events := make([]struct{t float64; idx int}, n)
   for i := 0; i < n; i++ {
       a[i] = float64(dist[i][u])
       b[i] = float64(dist[i][v])
       events[i] = struct{t float64; idx int}{(b[i] + bw - a[i]) / 2.0, i}
   }
   sort.Slice(events, func(i, j int) bool { return events[i].t < events[j].t })
   moved := make([]bool, n)
   var alpha, beta MaxHeap
   alpha = make(MaxHeap, 0, n)
   for i := 0; i < n; i++ {
       alpha = append(alpha, Pair{a[i], i})
   }
   heap.Init(&alpha)
   heap.Init(&beta)
   k := 0
   for k < n && events[k].t <= 0 {
       i := events[k].idx
       moved[i] = true
       heap.Push(&beta, Pair{b[i] + bw, i})
       k++
   }
   prevT := 0.0
   best := math.Inf(1)
   process := func(L, R float64) {
       if L > bw || R < 0 {
           return
       }
       l := L
       if l < 0 {
           l = 0
       }
       r := R
       if r > bw {
           r = bw
       }
       if l > r {
           return
       }
       for alpha.Len() > 0 {
           top := alpha[0]
           if moved[top.idx] {
               heap.Pop(&alpha)
               continue
           }
           break
       }
       var aMax float64
       if alpha.Len() > 0 {
           aMax = alpha[0].val
       } else {
           aMax = -1e30
       }
       var bMax float64
       if beta.Len() > 0 {
           bMax = beta[0].val
       } else {
           bMax = -1e30
       }
       t := (bMax - aMax) / 2.0
       if t < l {
           t = l
       } else if t > r {
           t = r
       }
       d1 := aMax + t
       d2 := bMax - t
       d := d1
       if d2 > d {
           d = d2
       }
       if d < best {
           best = d
       }
   }
   for k < n {
       currT := events[k].t
       process(prevT, currT)
       for k < n && events[k].t == currT {
           i := events[k].idx
           moved[i] = true
           heap.Push(&beta, Pair{b[i] + bw, i})
           k++
       }
       prevT = currT
   }
   process(prevT, bw)
   return best
}
