package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU (Disjoint Set Union) structure
type DSU struct {
   parent []int
   size   []int64
}

func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent: parent, size: size}
}

func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return
   }
   // union by size
   if d.size[rx] < d.size[ry] {
       rx, ry = ry, rx
   }
   d.parent[ry] = rx
   d.size[rx] += d.size[ry]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   gov := make([]int, k)
   for i := 0; i < k; i++ {
       fmt.Fscan(reader, &gov[i])
   }
   dsu := NewDSU(n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       dsu.Union(u, v)
   }
   // mark government roots and collect their sizes
   govRoot := make(map[int]bool)
   govSizes := make([]int64, 0, k)
   for _, c := range gov {
       r := dsu.Find(c)
       if !govRoot[r] {
           govRoot[r] = true
           govSizes = append(govSizes, dsu.size[r])
       }
   }
   // sum sizes of components without government
   var extra int64
   for i := 1; i <= n; i++ {
       if dsu.parent[i] == i {
           if !govRoot[i] {
               extra += dsu.size[i]
           }
       }
   }
   // find largest government component
   var maxIdx int
   for i := range govSizes {
       if i == 0 || govSizes[i] > govSizes[maxIdx] {
           maxIdx = i
       }
   }
   // compute total possible edges
   var totalEdges int64
   for i, sz := range govSizes {
       if i == maxIdx {
           s := sz + extra
           totalEdges += s * (s - 1) / 2
       } else {
           totalEdges += sz * (sz - 1) / 2
       }
   }
   // answer is totalEdges - existing m
   ans := totalEdges - int64(m)
   fmt.Println(ans)
}
