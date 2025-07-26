package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = int64(4e18)

// Edge for Dinic
type Edge struct { to, rev int; cap int64 }

// Dinic struct
type Dinic struct { N int; G [][]Edge; level []int; ptr []int }

func NewDinic(n int) *Dinic {
   d := &Dinic{N: n, G: make([][]Edge, n), level: make([]int, n), ptr: make([]int, n)}
   return d
}

func (d *Dinic) AddEdge(u, v int, c int64) {
   d.G[u] = append(d.G[u], Edge{v, len(d.G[v]), c})
   d.G[v] = append(d.G[v], Edge{u, len(d.G[u]) - 1, 0})
}

func (d *Dinic) bfs(s, t int) bool {
   for i := range d.level {
       d.level[i] = -1
   }
   q := make([]int, 0, d.N)
   d.level[s] = 0
   q = append(q, s)
   for i := 0; i < len(q); i++ {
       u := q[i]
       for _, e := range d.G[u] {
           if d.level[e.to] < 0 && e.cap > 0 {
               d.level[e.to] = d.level[u] + 1
               q = append(q, e.to)
           }
       }
   }
   return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t int, f int64) int64 {
   if u == t || f == 0 {
       return f
   }
   for &p := &d.ptr[u]; p < len(d.G[u]); p++ {
       e := &d.G[u][p]
       if d.level[e.to] == d.level[u]+1 && e.cap > 0 {
           pushed := d.dfs(e.to, t, min64(f, e.cap))
           if pushed > 0 {
               e.cap -= pushed
               d.G[e.to][e.rev].cap += pushed
               return pushed
           }
       }
   }
   return 0
}

func (d *Dinic) MaxFlow(s, t int) int64 {
   flow := int64(0)
   for d.bfs(s, t) {
       for i := range d.ptr {
           d.ptr[i] = 0
       }
       for {
           pushed := d.dfs(s, t, INF)
           if pushed == 0 {
               break
           }
           flow += pushed
       }
   }
   return flow
}

func min64(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m int
   fmt.Fscan(in, &n, &m)
   rects := make([][4]int, m)
   xs := make([]int, 0, 2*m)
   ys := make([]int, 0, 2*m)
   for i := 0; i < m; i++ {
       x1, y1, x2, y2 := 0, 0, 0, 0
       fmt.Fscan(in, &x1, &y1, &x2, &y2)
       rects[i] = [4]int{x1, y1, x2, y2}
       xs = append(xs, x1)
       xs = append(xs, x2+1)
       ys = append(ys, y1)
       ys = append(ys, y2+1)
   }
   if m == 0 {
       fmt.Println(0)
       return
   }
   sort.Ints(xs)
   xs = unique(xs)
   sort.Ints(ys)
   ys = unique(ys)
   nx := len(xs) - 1
   ny := len(ys) - 1
   // interval lengths
   rowLen := make([]int64, nx)
   for i := 0; i < nx; i++ {
       rowLen[i] = int64(xs[i+1] - xs[i])
   }
   colLen := make([]int64, ny)
   for j := 0; j < ny; j++ {
       colLen[j] = int64(ys[j+1] - ys[j])
   }
   // hasBlack
   has := make([][]bool, nx)
   for i := range has {
       has[i] = make([]bool, ny)
   }
   // mark
   for _, r := range rects {
       x1, y1, x2, y2 := r[0], r[1], r[2], r[3]
       xi1 := sort.SearchInts(xs, x1)
       xi2 := sort.SearchInts(xs, x2+1)
       yi1 := sort.SearchInts(ys, y1)
       yi2 := sort.SearchInts(ys, y2+1)
       for i := xi1; i < xi2; i++ {
           for j := yi1; j < yi2; j++ {
               has[i][j] = true
           }
       }
   }
   // build flow
   // nodes: 0=src, 1..nx rows, nx+1..nx+ny cols, nx+ny+1 sink
   N := nx + ny + 2
   src := 0
   sink := N - 1
   d := NewDinic(N)
   // source to rows
   for i := 0; i < nx; i++ {
       if rowLen[i] > 0 {
           d.AddEdge(src, 1+i, rowLen[i])
       }
   }
   // cols to sink
   for j := 0; j < ny; j++ {
       if colLen[j] > 0 {
           d.AddEdge(1+nx+j, sink, colLen[j])
       }
   }
   // edges
   for i := 0; i < nx; i++ {
       for j := 0; j < ny; j++ {
           if has[i][j] {
               d.AddEdge(1+i, 1+nx+j, INF)
           }
       }
   }
   flow := d.MaxFlow(src, sink)
   fmt.Println(flow)
}

// unique returns unique sorted slice
func unique(a []int) []int {
   res := a[:1]
   for i := 1; i < len(a); i++ {
       if a[i] != a[i-1] {
           res = append(res, a[i])
       }
   }
   return res
}
