package main

import (
   "bufio"
   "container/heap"
   "fmt"
   "math/big"
   "os"
)

const MOD = 1000000007

// Edge represents a graph edge with power exponent x (weight = 2^x)
type Edge struct {
   to, x int
}

// Item for priority queue
type Item struct {
   node int
   dist *big.Int
}

// Priority queue of Items
type PQ []*Item

func (pq PQ) Len() int { return len(pq) }
func (pq PQ) Less(i, j int) bool {
   // smaller dist has higher priority
   return pq[i].dist.Cmp(pq[j].dist) < 0
}
func (pq PQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PQ) Pop() interface{} {
   old := *pq
   n := len(old)
   it := old[n-1]
   *pq = old[0 : n-1]
   return it
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   adj := make([][]Edge, n+1)
   for i := 0; i < m; i++ {
       var u, v, x int
       fmt.Fscan(reader, &u, &v, &x)
       adj[u] = append(adj[u], Edge{v, x})
       adj[v] = append(adj[v], Edge{u, x})
   }
   var s, t int
   fmt.Fscan(reader, &s, &t)

   // Dijkstra with big.Int distances
   dist := make([]*big.Int, n+1)
   parent := make([]int, n+1)
   for i := 1; i <= n; i++ {
       dist[i] = nil
       parent[i] = -1
   }
   dist[s] = big.NewInt(0)

   pq := &PQ{}
   heap.Init(pq)
   heap.Push(pq, &Item{node: s, dist: big.NewInt(0)})

   for pq.Len() > 0 {
       it := heap.Pop(pq).(*Item)
       u := it.node
       d := it.dist
       if dist[u] == nil || d.Cmp(dist[u]) != 0 {
           continue
       }
       if u == t {
           break
       }
       for _, e := range adj[u] {
           // compute new distance = d + 2^e.x
           w := new(big.Int).Lsh(big.NewInt(1), uint(e.x))
           nd := new(big.Int).Add(d, w)
           if dist[e.to] == nil || nd.Cmp(dist[e.to]) < 0 {
               dist[e.to] = nd
               parent[e.to] = u
               heap.Push(pq, &Item{node: e.to, dist: nd})
           }
       }
   }

   if dist[t] == nil {
       fmt.Fprintln(writer, -1)
       return
   }
   // output distance mod
   modVal := new(big.Int).Mod(dist[t], big.NewInt(MOD)).Int64()
   fmt.Fprintln(writer, modVal)
   // reconstruct path
   path := []int{}
   for v := t; v != -1; v = parent[v] {
       path = append(path, v)
   }
   // reverse
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   fmt.Fprintln(writer, len(path))
   for i, v := range path {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
