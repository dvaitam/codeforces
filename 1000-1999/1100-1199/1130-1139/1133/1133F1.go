package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU (Disjoint Set Union) for union-find operations
type DSU struct {
   parent []int
}

// NewDSU initializes a DSU for elements 1..n
func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
   }
   return &DSU{parent: p}
}

// Find returns the representative of x with path compression
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union merges the sets containing x and y
func (d *DSU) Union(x, y int) {
   fx := d.Find(x)
   fy := d.Find(y)
   if fx != fy {
       d.parent[fx] = fy
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   u := make([]int, m)
   v := make([]int, m)
   deg := make([]int, n+1)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &u[i], &v[i])
       deg[u[i]]++
       deg[v[i]]++
   }
   // find vertex with maximum degree
   who := 1
   for i := 2; i <= n; i++ {
       if deg[i] > deg[who] {
           who = i
       }
   }
   dsu := NewDSU(n)
   used := make([]bool, m)
   // first, connect all edges incident to 'who'
   for i := 0; i < m; i++ {
       if u[i] == who || v[i] == who {
           dsu.Union(u[i], v[i])
           used[i] = true
       }
   }
   // then, add remaining edges to form spanning tree
   for i := 0; i < m; i++ {
       if used[i] {
           continue
       }
       if dsu.Find(u[i]) != dsu.Find(v[i]) {
           dsu.Union(u[i], v[i])
           used[i] = true
       }
   }
   // output the spanning tree edges
   for i := 0; i < m; i++ {
       if used[i] {
           fmt.Fprintln(writer, u[i], v[i])
       }
   }
}
