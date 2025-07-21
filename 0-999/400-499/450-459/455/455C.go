package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU with diameter maintenance
type DSU struct {
   parent []int
   size   []int
   diam   []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   sz := make([]int, n+1)
   d := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
       sz[i] = 1
       d[i] = 0
   }
   return &DSU{parent: p, size: sz, diam: d}
}

func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union by size, returns new root
func (d *DSU) Union(x, y int, newDia int) int {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return rx
   }
   // attach smaller to larger
   if d.size[rx] < d.size[ry] {
       rx, ry = ry, rx
   }
   d.parent[ry] = rx
   d.size[rx] += d.size[ry]
   // update diameter
   if newDia > d.diam[rx] {
       d.diam[rx] = newDia
   }
   return rx
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m, q int
   fmt.Fscan(in, &n, &m, &q)
   adj := make([][]int, n+1)
   dsu := NewDSU(n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       // initial union
       dsu.Union(u, v, 0)
   }
   // compute initial diameters per component
   visited := make([]bool, n+1)
   visited2 := make([]bool, n+1)
   dist1 := make([]int, n+1)
   depth2 := make([]int, n+1)
   var queue []int
   var compNodes []int
   for i := 1; i <= n; i++ {
       if visited[i] {
           continue
       }
      // BFS1: from i, find farthest u, collect comp nodes
       queue = queue[:0]
       compNodes = compNodes[:0]
       queue = append(queue, i)
       visited[i] = true
       dist1[i] = 0
       compNodes = append(compNodes, i)
       var u, maxd int
       for head := 0; head < len(queue); head++ {
           v := queue[head]
           d0 := dist1[v]
           if d0 > maxd {
               maxd = d0
               u = v
           }
           for _, w := range adj[v] {
               if !visited[w] {
                   visited[w] = true
                   dist1[w] = d0 + 1
                   queue = append(queue, w)
                   compNodes = append(compNodes, w)
               }
           }
       }
       // BFS2: from u, find diameter
       for _, v := range compNodes {
           visited2[v] = false
       }
       queue = queue[:0]
       queue = append(queue, u)
       visited2[u] = true
       depth2[u] = 0
       var dia int
       for head := 0; head < len(queue); head++ {
           v := queue[head]
           d0 := depth2[v]
           if d0 > dia {
               dia = d0
           }
           for _, w := range adj[v] {
               if !visited2[w] {
                   visited2[w] = true
                   depth2[w] = d0 + 1
                   queue = append(queue, w)
               }
           }
       }
       root := dsu.Find(i)
       dsu.diam[root] = dia
   }
   // process queries
   for j := 0; j < q; j++ {
       var t, x, y int
       fmt.Fscan(in, &t, &x)
       if t == 1 {
           r := dsu.Find(x)
           fmt.Fprintln(out, dsu.diam[r])
       } else {
           fmt.Fscan(in, &y)
           rx := dsu.Find(x)
           ry := dsu.Find(y)
           if rx == ry {
               continue
           }
           d1 := dsu.diam[rx]
           d2 := dsu.diam[ry]
           r1 := (d1 + 1) / 2
           r2 := (d2 + 1) / 2
           newD := d1
           if d2 > newD {
               newD = d2
           }
           if r1 + r2 + 1 > newD {
               newD = r1 + r2 + 1
           }
           dsu.Union(rx, ry, newD)
       }
   }
}
