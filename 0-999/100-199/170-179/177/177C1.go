package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU (Disjoint Set Union) for union-find
type DSU struct {
   parent []int
   size   []int
}

// NewDSU creates a new DSU with n elements (1-indexed)
func NewDSU(n int) *DSU {
   parent := make([]int, n+1)
   size := make([]int, n+1)
   for i := 1; i <= n; i++ {
       parent[i] = i
       size[i] = 1
   }
   return &DSU{parent: parent, size: size}
}

// Find returns the representative of x
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union merges the sets of x and y
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
   var n, k int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   fmt.Fscan(reader, &k)
   dsu := NewDSU(n)
   for i := 0; i < k; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       dsu.Union(u, v)
   }
   var m int
   fmt.Fscan(reader, &m)
   // invalid component roots due to dislike inside component
   invalid := make([]bool, n+1)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       ru := dsu.Find(u)
       rv := dsu.Find(v)
       if ru == rv {
           invalid[ru] = true
       }
   }
   // track max size among valid components
   maxSize := 0
   // To avoid recounting same root, iterate 1..n and check roots
   seen := make([]bool, n+1)
   for i := 1; i <= n; i++ {
       r := dsu.Find(i)
       if seen[r] {
           continue
       }
       seen[r] = true
       if invalid[r] {
           continue
       }
       if dsu.size[r] > maxSize {
           maxSize = dsu.size[r]
       }
   }
   fmt.Println(maxSize)
}
