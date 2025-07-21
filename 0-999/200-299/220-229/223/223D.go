package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math"
   "os"
   "sort"
)

type Edge struct {
   to int
   w  float64
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   xs := make([]int, n)
   ys := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &xs[i], &ys[i])
   }
   var s, t int
   fmt.Fscan(reader, &s, &t)
   s--
   t--
   // build graph
   N := n
   adj := make([][]Edge, N)
   // boundary edges
   for i := 0; i < n; i++ {
       j := (i + 1) % n
       dx := float64(xs[i] - xs[j])
       dy := float64(ys[i] - ys[j])
       d := math.Hypot(dx, dy)
       adj[i] = append(adj[i], Edge{j, d})
       adj[j] = append(adj[j], Edge{i, d})
   }
   // vertical descents between vertices with same x
   mp := make(map[int][]int)
   for i := 0; i < n; i++ {
       x := xs[i]
       mp[x] = append(mp[x], i)
   }
   for _, vec := range mp {
       sort.Slice(vec, func(i, j int) bool {
           return ys[vec[i]] > ys[vec[j]]
       })
       for k := 0; k+1 < len(vec); k++ {
           u := vec[k]
           v := vec[k+1]
           dy := float64(ys[u] - ys[v])
           if dy >= 0 {
               adj[u] = append(adj[u], Edge{v, dy})
           }
       }
   }
   // Dijkstra
   const inf = 1e300
   dist := make([]float64, N)
   for i := range dist {
       dist[i] = inf
   }
   dist[s] = 0
   visited := make([]bool, N)
   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, &Item{node: s, dist: 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(*Item)
       u := it.node
       if visited[u] {
           continue
       }
       visited[u] = true
       if u == t {
           break
       }
       for _, e := range adj[u] {
           v := e.to
           w := e.w
           nd := dist[u] + w
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(pq, &Item{node: v, dist: nd})
           }
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintf(out, "%.6f\n", dist[t])
}

// priority queue
type Item struct {
   node int
   dist float64
   idx  int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) {
   pq[i], pq[j] = pq[j], pq[i]
   pq[i].idx = i
   pq[j].idx = j
}
func (pq *PriorityQueue) Push(x interface{}) {
   it := x.(*Item)
   it.idx = len(*pq)
   *pq = append(*pq, it)
}
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   it.idx = -1
   *pq = old[:n-1]
   return it
}
