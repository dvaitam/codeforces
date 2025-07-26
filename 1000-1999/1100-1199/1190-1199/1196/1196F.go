package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

// Edge represents an adjacency list edge
type Edge struct {
   to   int
   w    int64
}

// Pair holds source vertex and distance
type Pair struct {
   src  int
   dist int64
}

// State for priority queue
type State struct {
   dist int64
   node int
   src  int
}

// PriorityQueue implements heap.Interface
type PriorityQueue []State

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(State)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   x := old[n-1]
   *pq = old[0 : n-1]
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, K int
   fmt.Fscan(reader, &n, &m, &K)
   adj := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var x, y int
       var w int64
       fmt.Fscan(reader, &x, &y, &w)
       adj[x] = append(adj[x], Edge{to: y, w: w})
       adj[y] = append(adj[y], Edge{to: x, w: w})
   }

   best := make([][]Pair, n+1)
   visited := make(map[int64]bool)
   ans := make([]int64, 0, K)

   pq := &PriorityQueue{}
   heap.Init(pq)
   for i := 1; i <= n; i++ {
       heap.Push(pq, State{dist: 0, node: i, src: i})
   }

   for pq.Len() > 0 && len(ans) < K {
       s := heap.Pop(pq).(State)
       d, v, src := s.dist, s.node, s.src
       // check existing source
       skip := false
       for _, p := range best[v] {
           if p.src == src {
               skip = true
               break
           }
       }
       if skip {
           continue
       }
       // prune if too large
       if len(best[v]) >= K && best[v][len(best[v])-1].dist <= d {
           continue
       }
       // insert into best[v]
       arr := best[v]
       pos := sort.Search(len(arr), func(i int) bool { return arr[i].dist > d })
       // insert at pos
       arr = append(arr, Pair{})
       copy(arr[pos+1:], arr[pos:])
       arr[pos] = Pair{src: src, dist: d}
       if len(arr) > K {
           arr = arr[:K]
       }
       best[v] = arr
       // generate candidate paths
       for _, p := range best[v] {
           if p.src == src {
               continue
           }
           u := src
           wsrc := p.src
           if u > wsrc {
               u, wsrc = wsrc, u
           }
           key := (int64(u) << 32) | int64(wsrc)
           if !visited[key] {
               visited[key] = true
               ans = append(ans, d + p.dist)
           }
       }
       // expand neighbors
       for _, e := range adj[v] {
           heap.Push(pq, State{dist: d + e.w, node: e.to, src: src})
       }
   }

   sort.Slice(ans, func(i, j int) bool { return ans[i] < ans[j] })
   if len(ans) >= K {
       fmt.Fprintln(writer, ans[K-1])
   }
}
