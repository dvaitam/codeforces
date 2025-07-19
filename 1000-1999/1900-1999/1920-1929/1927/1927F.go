package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU (Disjoint Set Union) with union by size and path compression
type DSU struct {
   parent, size []int
}

// NewDSU initializes a DSU for 1..n
func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int, n+1)
   for i := 0; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent: parent, size: size}
}

// Find returns the representative of u
func (d *DSU) Find(u int) int {
   if d.parent[u] != u {
       d.parent[u] = d.Find(d.parent[u])
   }
   return d.parent[u]
}

// Union merges the sets of u and v
func (d *DSU) Union(u, v int) {
   ru := d.Find(u)
   rv := d.Find(v)
   if ru == rv {
       return
   }
   if d.size[ru] < d.size[rv] {
       d.parent[ru] = rv
       d.size[rv] += d.size[ru]
   } else {
       d.parent[rv] = ru
       d.size[ru] += d.size[rv]
   }
}

// Edge represents an undirected edge with weight
type Edge struct {
   u, v, w int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       edges := make([]Edge, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &edges[i].u, &edges[i].v, &edges[i].w)
       }
       // sort edges by descending weight
       sort.Slice(edges, func(i, j int) bool {
           return edges[i].w > edges[j].w
       })
       dsu := NewDSU(n)
       // adjacency list for graph built so far
       g := make([][]int, n+1)
       var st, en, cost int
       for _, e := range edges {
           u, v, w := e.u, e.v, e.w
           // add edge to graph
           g[u] = append(g[u], v)
           g[v] = append(g[v], u)
           // if u and v already connected, this edge closes a cycle
           if dsu.Find(u) == dsu.Find(v) {
               st, en, cost = u, v, w
           }
           dsu.Union(u, v)
       }
       // BFS to find path from st to en excluding the direct closing edge
       vis := make([]bool, n+1)
       parent := make([]int, n+1)
       queue := make([]int, 0, n)
       queue = append(queue, st)
       vis[st] = true
       // bfs
       for head := 0; head < len(queue); head++ {
           u := queue[head]
           if u == en {
               break
           }
           for _, v2 := range g[u] {
               if u == st && v2 == en {
                   // skip the direct cycle edge
                   continue
               }
               if !vis[v2] {
                   vis[v2] = true
                   parent[v2] = u
                   queue = append(queue, v2)
               }
           }
       }
       // reconstruct path from en to st
       path := make([]int, 0)
       cur := en
       for {
           path = append(path, cur)
           if cur == st {
               break
           }
           cur = parent[cur]
       }
       // output
       fmt.Fprintln(writer, cost, len(path))
       for _, x := range path {
           fmt.Fprint(writer, x, " ")
       }
       fmt.Fprintln(writer)
   }
}
