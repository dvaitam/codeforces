package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

// Edge represents a road between u and v with weight w
type Edge struct {
   u, v int
   w    int64
}

// Item is a node in the priority queue
type Item struct {
   v    int
   dist int64
}

// MinHeap implements a min-heap of Items
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
   item := old[n-1]
   *h = old[0 : n-1]
   return item
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, s int
   var L int64
   if _, err := fmt.Fscan(reader, &n, &m, &s); err != nil {
       return
   }
   adj := make([][]struct{to int; w int64}, n+1)
   edges := make([]Edge, m)
   for i := 0; i < m; i++ {
       var u, v int
       var w int64
       fmt.Fscan(reader, &u, &v, &w)
       adj[u] = append(adj[u], struct{to int; w int64}{v, w})
       adj[v] = append(adj[v], struct{to int; w int64}{u, w})
       edges[i] = Edge{u: u, v: v, w: w}
   }
   fmt.Fscan(reader, &L)
   const INF = int64(4e18)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   // Dijkstra from s
   pq := &MinHeap{}
   heap.Init(pq)
   dist[s] = 0
   heap.Push(pq, Item{v: s, dist: 0})
   for pq.Len() > 0 {
       cur := heap.Pop(pq).(Item)
       if cur.dist != dist[cur.v] {
           continue
       }
       for _, e := range adj[cur.v] {
           nd := cur.dist + e.w
           if nd < dist[e.to] {
               dist[e.to] = nd
               heap.Push(pq, Item{v: e.to, dist: nd})
           }
       }
   }
   // count silos at cities
   var ans int64
   for i := 1; i <= n; i++ {
       if dist[i] == L {
           ans++
       }
   }
   // check silos on roads
   for _, e := range edges {
       du := dist[e.u]
       dv := dist[e.v]
       w := e.w
       // candidate from u
       flagU := du < L && du+w > L && du+dv+w >= 2*L
       // candidate from v
       flagV := dv < L && dv+w > L && du+dv+w >= 2*L
       if flagU && flagV && du+dv+w == 2*L {
           ans++
       } else {
           if flagU {
               ans++
           }
           if flagV {
               ans++
           }
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
