package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU implements disjoint set union with union by size and path compression.
type DSU struct {
   p []int
}

// NewDSU creates a DSU for n elements (0..n-1).
func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := range p {
       p[i] = -1
   }
   return &DSU{p: p}
}

// Find finds the representative of x.
func (d *DSU) Find(x int) int {
   if d.p[x] < 0 {
       return x
   }
   d.p[x] = d.Find(d.p[x])
   return d.p[x]
}

// Union merges the sets of a and b. Returns true if merged, false if already same set.
func (d *DSU) Union(a, b int) bool {
   a = d.Find(a)
   b = d.Find(b)
   if a == b {
       return false
   }
   // Union by size (p stores negative size)
   if d.p[a] > d.p[b] {
       a, b = b, a
   }
   d.p[a] += d.p[b]
   d.p[b] = a
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var m, n int
   fmt.Fscan(reader, &m, &n)
   a := make([]int64, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, n)
   for j := 0; j < n; j++ {
       fmt.Fscan(reader, &b[j])
   }
   type Edge struct{ u, v int; w int64 }
   edges := make([]Edge, 0)
   var totalCost int64
   // Read sets and prepare edges in bipartite graph between sets and elements
   for i := 0; i < m; i++ {
       var sz int
       fmt.Fscan(reader, &sz)
       for k := 0; k < sz; k++ {
           var j int
           fmt.Fscan(reader, &j)
           j--
           w := a[i] + b[j]
           edges = append(edges, Edge{u: i, v: m + j, w: w})
           totalCost += w
       }
   }
   // Maximum spanning forest: keep highest weights without forming cycles
   sort.Slice(edges, func(i, j int) bool { return edges[i].w > edges[j].w })
   dsu := NewDSU(m + n)
   var keepSum int64
   for _, e := range edges {
       if dsu.Union(e.u, e.v) {
           keepSum += e.w
       }
   }
   // Deletion cost is totalCost - sum of kept weights
   res := totalCost - keepSum
   fmt.Fprintln(writer, res)
}
