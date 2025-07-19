package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "os"
)

const inf = int64(1e18)
const cFijo = int64(1000000)

// Edge represents a neighbor and its status z (1=good, 0=broken)
type Edge struct {
   to int
   z  int
}

// Item for priority queue
type Item struct {
   node int
   dist int64
}

// PriorityQueue implements heap.Interface
type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[:n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M int
   fmt.Fscan(reader, &N, &M)
   g := make([][]Edge, N)
   type EdgeRec struct{ u, v, z int }
   edges := make([]EdgeRec, 0, M)
   rotos := 0
   for i := 0; i < M; i++ {
       var u, v, z int
       fmt.Fscan(reader, &u, &v, &z)
       u--
       v--
       g[u] = append(g[u], Edge{v, z})
       g[v] = append(g[v], Edge{u, z})
       edges = append(edges, EdgeRec{u, v, z})
       if z == 0 {
           rotos++
       }
   }
   // Dijkstra
   dist := make([]int64, N)
   parent := make([]int, N)
   for i := 0; i < N; i++ {
       dist[i] = inf
       parent[i] = -1
   }
   dist[0] = 0
   pq := &PriorityQueue{}
   heap.Init(pq)
   heap.Push(pq, Item{node: 0, dist: 0})
   for pq.Len() > 0 {
       it := heap.Pop(pq).(Item)
       u := it.node
       d := it.dist
       if d != dist[u] {
           continue
       }
       for _, e := range g[u] {
           cost := cFijo
           if e.z == 0 {
               cost++
           }
           nd := d + cost
           v := e.to
           if nd < dist[v] {
               dist[v] = nd
               parent[v] = u
               heap.Push(pq, Item{node: v, dist: nd})
           }
       }
   }
   // Reconstruct path
   pathMap := make(map[int64]bool)
   // map of edge key to z
   zMap := make(map[int64]int)
   for _, e := range edges {
       u, v := e.u, e.v
       if u > v {
           u, v = v, u
       }
       key := int64(u)<<32 | int64(v)
       zMap[key] = e.z
   }
   brokenInPath := 0
   pathLen := 0
   cur := N - 1
   for parent[cur] != -1 {
       p := parent[cur]
       u, v := p, cur
       if u > v {
           u, v = v, u
       }
       key := int64(u)<<32 | int64(v)
       pathMap[key] = true
       if zMap[key] == 0 {
           brokenInPath++
       }
       pathLen++
       cur = p
   }
   // Compute operations
   ops := M - rotos - pathLen + 2*brokenInPath
   fmt.Fprintln(writer, ops)
   // Emit flips
   for _, e := range edges {
       u, v, z := e.u, e.v, e.z
       uu, vv := u, v
       if uu > vv {
           uu, vv = vv, uu
       }
       key := int64(uu)<<32 | int64(vv)
       if pathMap[key] {
           // on path: broken-> make good
           if z == 0 {
               fmt.Fprintf(writer, "%d %d 1\n", u+1, v+1)
           }
       } else {
           // off path: good-> make broken
           if z == 1 {
               fmt.Fprintf(writer, "%d %d 0\n", u+1, v+1)
           }
       }
   }
}
