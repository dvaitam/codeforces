package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Disjoint set union structure
type DSU struct {
   parent []int
   size   []int64
}

// NewDSU initializes DSU for n elements
func NewDSU(n int) *DSU {
   p := make([]int, n)
   sz := make([]int64, n)
   for i := 0; i < n; i++ {
       p[i] = i
       sz[i] = 1
   }
   return &DSU{parent: p, size: sz}
}

// Find with path compression
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// UnionByWeight merges x and y, returns added pairs count * weight if merged
func (d *DSU) Union(x, y int, weight int) int64 {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return 0
   }
   // union by size
   if d.size[rx] < d.size[ry] {
       rx, ry = ry, rx
   }
   // all pairs between ry and rx have min node weight = weight
   contrib := int64(weight) * d.size[rx] * d.size[ry]
   d.parent[ry] = rx
   d.size[rx] += d.size[ry]
   return contrib
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // build adjacency list
   g := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
       g[v] = append(g[v], u)
   }
   // order nodes by descending weight
   ord := make([]int, n)
   for i := 0; i < n; i++ {
       ord[i] = i
   }
   sort.Slice(ord, func(i, j int) bool {
       return a[ord[i]] > a[ord[j]]
   })
   dsu := NewDSU(n)
   active := make([]bool, n)
   var sum int64
   // activate nodes and union with active neighbors
   for _, u := range ord {
       active[u] = true
       for _, v := range g[u] {
           if active[v] {
               sum += dsu.Union(u, v, a[u])
           }
       }
   }
   // compute average over ordered pairs p!=q
   denom := float64(n) * float64(n-1)
   avg := 2.0 * float64(sum) / denom
   fmt.Fprintf(out, "%.6f\n", avg)
}
