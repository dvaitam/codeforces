package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
)

const LG = 19
const INF = math.MaxInt64 / 4

type Edge struct {
   to    int
   c1, c2 int64
}

type Item struct {
   dist int64
   node int
}
// A MinHeap implements heap.Interface for Items.
type MinHeap []Item
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
   *h = append(*h, x.(Item))
}
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   x := old[n-1]
   *h = old[0 : n-1]
   return x
}

type mat [2][2]int64

func merge(a, b mat) mat {
   var c mat
   // c[i][j] = min(a[i][0]+b[0][j], a[i][1]+b[1][j])
   for i := 0; i < 2; i++ {
       for j := 0; j < 2; j++ {
           x := a[i][0] + b[0][j]
           y := a[i][1] + b[1][j]
           if x < y {
               c[i][j] = x
           } else {
               c[i][j] = y
           }
       }
   }
   return c
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   w := make([]int64, n+1)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &w[i])
       dist[i] = w[i]
   }
   v := make([][]Edge, n+1)
   for i := 1; i < n; i++ {
       var u, vv int
       var c1, c2 int64
       fmt.Fscan(in, &u, &vv, &c1, &c2)
       v[u] = append(v[u], Edge{vv, c1, c2})
       v[vv] = append(v[vv], Edge{u, c1, c2})
   }
   // Dijkstra
   h := &MinHeap{}
   heap.Init(h)
   for i := 1; i <= n; i++ {
       heap.Push(h, Item{dist[i], i})
   }
   for h.Len() > 0 {
       it := heap.Pop(h).(Item)
       d := it.dist
       u := it.node
       if d != dist[u] {
           continue
       }
       for _, e := range v[u] {
           nd := d + e.c1 + e.c2
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(h, Item{nd, e.to})
           }
       }
   }
   w = dist
   // Prepare LCA and matrices
   par := make([][LG]int, n+1)
   depth := make([]int, n+1)
   to := make([][LG]mat, n+1)
   // Initialize to[][] to INF
   for u := 1; u <= n; u++ {
       for i := 0; i < LG; i++ {
           for x := 0; x < 2; x++ {
               for y := 0; y < 2; y++ {
                   to[u][i][x][y] = INF
               }
           }
       }
   }
   // BFS from 1
   queue := make([]int, 0, n)
   queue = append(queue, 1)
   depth[1] = 0
   par[1][0] = 0
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       // compute binary lifts for u
       for i := 1; i < LG; i++ {
           pu := par[u][i-1]
           par[u][i] = par[pu][i-1]
           to[u][i] = merge(to[u][i-1], to[pu][i-1])
       }
       for _, e := range v[u] {
           v2 := e.to
           if v2 == par[u][0] {
               continue
           }
           par[v2][0] = u
           depth[v2] = depth[u] + 1
           // initialize to[v2][0]
           to[v2][0][0][0] = min64(e.c1, w[u]+w[v2]+e.c2)
           to[v2][0][1][1] = min64(e.c2, w[u]+w[v2]+e.c1)
           to[v2][0][0][1] = min64(e.c1+w[u], w[v2]+e.c2)
           to[v2][0][1][0] = min64(e.c2+w[u], w[v2]+e.c1)
           queue = append(queue, v2)
       }
   }
   // queries
   var q int
   fmt.Fscan(in, &q)
   var ta, tb int
   var aa, bb int
   var ans int64
   at := mat{}
   bt := mat{}
   var ct mat
   for q > 0 {
       q--
       fmt.Fscan(in, &ta, &tb)
       if ta&1 == 1 {
           aa = 0
           ta = (ta + 1) / 2
       } else {
           aa = 1
           ta = ta / 2
       }
       if tb&1 == 1 {
           bb = 0
           tb = (tb + 1) / 2
       } else {
           bb = 1
           tb = tb / 2
       }
       // reset at, bt to zero
       for x := 0; x < 2; x++ {
           for y := 0; y < 2; y++ {
               at[x][y] = 0
               bt[x][y] = 0
           }
       }
       at[0][1], at[1][0] = w[ta], w[ta]
       bt[0][1], bt[1][0] = w[tb], w[tb]
       // lift tb to depth ta
       if depth[ta] > depth[tb] {
           ta, tb = tb, ta
           aa, bb = bb, aa
       }
       diff := depth[tb] - depth[ta]
       for i := 0; i < LG; i++ {
           if diff>>i&1 == 1 {
               ct = merge(bt, to[tb][i])
               bt = ct
               tb = par[tb][i]
           }
       }
       if ta != tb {
           for i := LG - 1; i >= 0; i-- {
               if par[ta][i] != par[tb][i] {
                   ct = merge(at, to[ta][i])
                   at = ct
                   ct = merge(bt, to[tb][i])
                   bt = ct
                   ta = par[ta][i]
                   tb = par[tb][i]
               }
           }
           // one more to reach LCA
           ct = merge(at, to[ta][0])
           at = ct
           ct = merge(bt, to[tb][0])
           bt = ct
           ta = par[ta][0]
           tb = par[tb][0]
       }
       // compute answer
       x0 := at[aa][0] + bt[bb][0]
       x1 := at[aa][1] + bt[bb][1]
       if x1 < x0 {
           ans = x1
       } else {
           ans = x0
       }
       fmt.Fprintln(out, ans)
   }
}
