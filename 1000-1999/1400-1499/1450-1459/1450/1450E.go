package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const inf = 1e9

type Edge struct {
   to, val int
}

var (
   n, m    int
   graph   [][]Edge
   col     []int
   disSPFA []int
   cnt     []int
   inq     []bool
)

func addEdge(u, v, w int) {
   graph[u] = append(graph[u], Edge{v, w})
}

// dfs for bipartiteness, color with ±1
func dfs(u, c int) bool {
   col[u] = c
   for _, e := range graph[u] {
       v := e.to
       if col[v] != 0 {
           if col[v] != -c {
               return false
           }
       } else {
           if !dfs(v, -c) {
               return false
           }
       }
   }
   return true
}

// SPFA from 1, detect negative cycle
func SPFA() bool {
   for i := 1; i <= n; i++ {
       disSPFA[i] = inf
       cnt[i] = 0
       inq[i] = false
   }
   disSPFA[1] = 0
   q := make([]int, 0, n)
   q = append(q, 1)
   inq[1] = true
   for len(q) > 0 {
       u := q[0]
       q = q[1:]
       inq[u] = false
       for _, e := range graph[u] {
           v := e.to
           if disSPFA[v] > disSPFA[u] + e.val {
               disSPFA[v] = disSPFA[u] + e.val
               cnt[v] = cnt[u] + 1
               if cnt[v] > n {
                   return true
               }
               if !inq[v] {
                   inq[v] = true
                   q = append(q, v)
               }
           }
       }
   }
   return false
}

// Dijkstra with early exit on negative dist
func solveFrom(src int) (int, []int, bool) {
   dist := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = inf
   }
   dist[src] = 0
   // min-heap
   h := &MinHeap{}
   heap.Init(h)
   heap.Push(h, Item{0, src})
   for h.Len() > 0 {
       it := heap.Pop(h).(Item)
       d, u := it.dist, it.node
       if d > dist[u] {
           continue
       }
       if d < 0 {
           return 0, nil, false
       }
       for _, e := range graph[u] {
           v := e.to
           nd := d + e.val
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(h, Item{nd, v})
           }
       }
   }
   mx := 0
   for i := 1; i <= n; i++ {
       if dist[i] > mx {
           mx = dist[i]
       }
   }
   // drop index 0
   res := make([]int, n)
   for i := 1; i <= n; i++ {
       res[i-1] = dist[i]
   }
   return mx, res, true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n, &m)
   graph = make([][]Edge, n+1)
   col = make([]int, n+1)
   disSPFA = make([]int, n+1)
   cnt = make([]int, n+1)
   inq = make([]bool, n+1)
   for i := 0; i < m; i++ {
       var x, y, z int
       fmt.Fscan(in, &x, &y, &z)
       addEdge(x, y, 1)
       if z != 0 {
           addEdge(y, x, -1)
       } else {
           addEdge(y, x, 1)
       }
   }
   // check bipartite
   if !dfs(1, 1) || SPFA() {
       fmt.Fprintln(out, "NO")
       return
   }
   best := -1
   var bestVec []int
   for i := 1; i <= n; i++ {
       mx, vec, ok := solveFrom(i)
       if ok && mx > best {
           best = mx
           bestVec = vec
       }
   }
   fmt.Fprintln(out, "YES")
   fmt.Fprintln(out, best)
   for i, v := range bestVec {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}

// priority queue
type Item struct {
   dist, node int
}

type MinHeap []Item

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h MinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
   old := *h
   n := len(old)
   item := old[n-1]
   *h = old[0 : n-1]
   return item
}
