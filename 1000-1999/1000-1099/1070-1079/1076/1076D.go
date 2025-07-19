package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

type Edge struct {
   to, w, idx int
}

type Item struct {
   node int
   dist int64
}

// Min-heap of Items
type PriorityQueue []Item

func (h PriorityQueue) Len() int            { return len(h) }
func (h PriorityQueue) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h PriorityQueue) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PriorityQueue) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *PriorityQueue) Pop() interface{} {
   old := *h
   n := len(old)
   item := old[n-1]
   *h = old[:n-1]
   return item
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var N, M, K int
   fmt.Fscan(in, &N, &M, &K)
   adj := make([][]Edge, N+1)
   for i := 1; i <= M; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       adj[a] = append(adj[a], Edge{to: b, w: c, idx: i})
       adj[b] = append(adj[b], Edge{to: a, w: c, idx: i})
   }

   const inf = int64(1e18)
   dist := make([]int64, N+1)
   for i := 1; i <= N; i++ {
       dist[i] = inf
   }
   dist[1] = 0
   parentEdge := make([]int, N+1)
   parentNode := make([]int, N+1)

   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{node: 1, dist: 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u := it.node
       if it.dist != dist[u] {
           continue
       }
       for _, e := range adj[u] {
           nd := it.dist + int64(e.w)
           if nd < dist[e.to] {
               dist[e.to] = nd
               parentEdge[e.to] = e.idx
               parentNode[e.to] = u
               heap.Push(pq, Item{node: e.to, dist: nd})
           }
       }
   }

   maxEdges := N - 1
   if K > maxEdges {
       K = maxEdges
   }
   // build tree of shortest path edges
   children := make([][]Edge, N+1)
   for v := 2; v <= N; v++ {
       p := parentNode[v]
       if p != 0 {
           children[p] = append(children[p], Edge{to: v, idx: parentEdge[v]})
       }
   }

   // DFS up to K edges
   ans := make([]int, 0, K)
   type stackEntry struct{ node, nxt int }
   stack := []stackEntry{{node: 1, nxt: 0}}
   for len(stack) > 0 && len(ans) < K {
       top := &stack[len(stack)-1]
       u := top.node
       if top.nxt < len(children[u]) {
           e := children[u][top.nxt]
           top.nxt++
           ans = append(ans, e.idx)
           stack = append(stack, stackEntry{node: e.to, nxt: 0})
       } else {
           stack = stack[:len(stack)-1]
       }
   }

   fmt.Fprintln(out, len(ans))
   for i, v := range ans {
       if i > 0 {
           fmt.Fprint(out, ' ')
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
