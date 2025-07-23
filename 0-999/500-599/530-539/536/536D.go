package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
   "sort"
)

// Dijkstra with min-heap
type edge struct { to int; w int64 }
func dijkstra(n int, src int, adj [][]edge) []int64 {
   dist := make([]int64, n)
   for i := range dist {
       dist[i] = math.MaxInt64
   }
   dist[src] = 0
   // min-heap of pairs (dist, node)
   h := &minHeap{}
   heap.Init(h)
   heap.Push(h, pair{0, src})
   for h.Len() > 0 {
       p := heap.Pop(h).(pair)
       d, u := p.d, p.u
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           nd := d + e.w
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(h, pair{nd, e.to})
           }
       }
   }
   return dist
}
type pair struct { d, u int64 }
type minHeap []pair
func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(pair)) }
func (h *minHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   var s, t int
   fmt.Fscan(in, &n, &m)
   fmt.Fscan(in, &s, &t)
   s--; t--
   p := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
   }
   adj := make([][]edge, n)
   for i := 0; i < m; i++ {
       var u, v int; var w int64
       fmt.Fscan(in, &u, &v, &w)
       u--; v--
       adj[u] = append(adj[u], edge{v, w})
       adj[v] = append(adj[v], edge{u, w})
   }
   ds := dijkstra(n, s, adj)
   dt := dijkstra(n, t, adj)
   // compress distances
   us := make([]int64, n)
   ut := make([]int64, n)
   copy(us, ds)
   copy(ut, dt)
   sort.Slice(us, func(i, j int) bool { return us[i] < us[j] })
   sort.Slice(ut, func(i, j int) bool { return ut[i] < ut[j] })
   us = unique(us)
   ut = unique(ut)
   ks := len(us)
   kt := len(ut)
   // map to index
   idxs := make([]int, n)
   idxt := make([]int, n)
   for i := 0; i < n; i++ {
       idxs[i] = sort.Search(len(us), func(j int) bool { return us[j] >= ds[i] }) + 1
       idxt[i] = sort.Search(len(ut), func(j int) bool { return ut[j] >= dt[i] }) + 1
   }
   // F[r][c]
   F := make([][]int64, ks+1)
   for i := range F {
       F[i] = make([]int64, kt+1)
   }
   for i := 0; i < n; i++ {
       r := idxs[i]
       c := idxt[i]
       F[r][c] += p[i]
   }
   // rowCumul and RP
   rowC := make([][]int64, ks+1)
   for i := 1; i <= ks; i++ {
       rowC[i] = make([]int64, kt+2)
       for j := kt; j >= 1; j-- {
           rowC[i][j] = rowC[i][j+1] + F[i][j]
       }
   }
   RP := make([][]int64, ks+1)
   for i := 0; i <= ks; i++ {
       RP[i] = make([]int64, kt+1)
   }
   for i := 1; i <= ks; i++ {
       for j := 0; j <= kt; j++ {
           RP[i][j] = RP[i-1][j] + rowC[i][j+1]
       }
   }
   // colCumul and CP
   colC := make([][]int64, ks+1)
   for i := 0; i <= ks; i++ {
       colC[i] = make([]int64, kt+1)
   }
   for c := 1; c <= kt; c++ {
       // sum for r>ks is 0
       for i := ks - 1; i >= 0; i-- {
           colC[i][c] = colC[i+1][c]
           if i+1 <= ks {
               colC[i][c] += F[i+1][c]
           }
       }
   }
   CP := make([][]int64, kt+1)
   for j := 0; j <= kt; j++ {
       CP[j] = make([]int64, ks+1)
   }
   for j := 1; j <= kt; j++ {
       for i := 0; i <= ks; i++ {
           CP[j][i] = CP[j-1][i] + colC[i][j]
       }
   }
   // dpT and dpNCol, mnTerm
   dpT := make([][]int64, ks+1)
   for i := range dpT {
       dpT[i] = make([]int64, kt+1)
   }
   dpNcol := make([]int64, ks+1)
   mnTerm := make([]int64, ks+1)
   const INF = math.MaxInt64 / 4
   // init for j=kt+? dpNcol[*]=0, mnTerm[*]=INF
   for i := 0; i <= ks; i++ {
       dpNcol[i] = 0
       mnTerm[i] = INF
   }
   // main DP
   mxT := make([]int64, ks+1)
   for j := kt; j >= 0; j-- {
       if j < kt {
           for i := 0; i <= ks; i++ {
               // dpN[i][j] = mnTerm[i]_(j+1) + CP[j][i]
               dpNcol[i] = mnTerm[i] + CP[j][i]
           }
       }
       // dpT col j
       // build mxT
       // start with i=ks
       mx := math.MinInt64
       for i := ks; i >= 0; i-- {
           v := RP[i][j] + dpNcol[i]
           if v > mx {
               mx = v
           }
           mxT[i] = mx
       }
       for i := 0; i <= ks; i++ {
           if i == ks {
               dpT[i][j] = 0
           } else {
               dpT[i][j] = mxT[i+1] - RP[i][j]
           }
       }
       // update mnTerm for next j
       for i := 0; i <= ks; i++ {
           v := dpT[i][j] - CP[j][i]
           if v < mnTerm[i] {
               mnTerm[i] = v
           }
       }
   }
   // result at state (0,0)
   res := dpT[0][0]
   if res > 0 {
       fmt.Fprintln(out, "Break a heart")
   } else if res < 0 {
       fmt.Fprintln(out, "Cry")
   } else {
       fmt.Fprintln(out, "Flowers")
   }
}

func unique(a []int64) []int64 {
   res := a[:0]
   prev := int64(math.MinInt64)
   for _, v := range a {
       if v != prev {
           res = append(res, v)
           prev = v
       }
   }
   return res
}
