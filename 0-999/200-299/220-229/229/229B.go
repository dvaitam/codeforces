package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
   "sort"
)

type edge struct {
   to int
   w  int64
}

// interval of forbidden departure times [l, r]
type interval struct{
   l, r int64
}

// priority queue item
type Item struct{
   node int
   dist int64
}

// min-heap of Items by dist
type PQ []Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[:n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]edge, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       var c int64
       fmt.Fscan(reader, &a, &b, &c)
       adj[a] = append(adj[a], edge{b, c})
       adj[b] = append(adj[b], edge{a, c})
   }
   forb := make([][]interval, n+1)
   for i := 1; i <= n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       if k == 0 {
           continue
       }
       times := make([]int64, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &times[j])
       }
       ivs := make([]interval, 0, k)
       start := times[0]
       prev := times[0]
       for j := 1; j < k; j++ {
           if times[j] == prev+1 {
               prev = times[j]
           } else {
               ivs = append(ivs, interval{start, prev})
               start = times[j]
               prev = times[j]
           }
       }
       ivs = append(ivs, interval{start, prev})
       forb[i] = ivs
   }

   const INF = int64(9e18)
   dist := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = INF
   }
   dist[1] = 0
   pq := &PQ{}
   heap.Init(pq)
   heap.Push(pq, Item{1, 0})

   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u := it.node
       d := it.dist
       if d != dist[u] {
           continue
       }
       if u == n {
           break
       }
       dep := d
       ivs := forb[u]
       if len(ivs) > 0 {
           idx := sort.Search(len(ivs), func(i int) bool {
               return ivs[i].l > d
           }) - 1
           if idx >= 0 && ivs[idx].l <= d && d <= ivs[idx].r {
               dep = ivs[idx].r + 1
           }
       }
       for _, e := range adj[u] {
           v := e.to
           nd := dep + e.w
           if nd < dist[v] {
               dist[v] = nd
               heap.Push(pq, Item{v, nd})
           }
       }
   }
   if dist[n] == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(dist[n])
   }
}
