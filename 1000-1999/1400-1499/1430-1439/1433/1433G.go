package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "io"
   "os"
)

type edge struct{ to, w int }

type Item struct {
   v    int
   dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

func dijkstra(n int, adj [][]edge, src int) []int64 {
   const INF = int64(4e18)
   dist := make([]int64, n)
   for i := range dist {
       dist[i] = INF
   }
   dist[src] = 0
   pq := &PriorityQueue{{v: src, dist: 0}}
   heap.Init(pq)
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u, d := it.v, it.dist
       if d != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           nd := d + int64(e.w)
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{v: e.to, dist: nd})
           }
       }
   }
   return dist
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err == io.EOF {
       return
   }
   adj := make([][]edge, n)
   edges := make([]struct{ u, v int }, m)
   for i := 0; i < m; i++ {
       var x, y, w int
       fmt.Fscan(in, &x, &y, &w)
       x--, y--
       adj[x] = append(adj[x], edge{to: y, w: w})
       adj[y] = append(adj[y], edge{to: x, w: w})
       edges[i] = struct{ u, v int }{u: x, v: y}
   }
   qs := make([]struct{ a, b int }, k)
   for i := 0; i < k; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       qs[i] = struct{ a, b int }{a: a - 1, b: b - 1}
   }
   // All-pairs shortest paths via Dijkstra from each node
   dist := make([][]int64, n)
   for i := 0; i < n; i++ {
       dist[i] = dijkstra(n, adj, i)
   }
   // Compute base sum
   var base int64
   for _, q := range qs {
       base += dist[q.a][q.b]
   }
   ans := base
   // Try zeroing each edge
   for i := 0; i < m; i++ {
       u := edges[i].u
       v := edges[i].v
       var sum int64
       for _, q := range qs {
           d0 := dist[q.a][q.b]
           d1 := dist[q.a][u] + dist[v][q.b]
           if d1 < d0 {
               d0 = d1
           }
           d2 := dist[q.a][v] + dist[u][q.b]
           if d2 < d0 {
               d0 = d2
           }
           sum += d0
           // early break if sum already >= ans
           if sum >= ans {
               break
           }
       }
       if sum < ans {
           ans = sum
       }
   }
   fmt.Println(ans)
}
