package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU (Disjoint Set Union) structure
type DSU struct {
   parent []int
   rank   []int
}

// NewDSU creates a DSU for n elements (0..n-1)
func NewDSU(n int) *DSU {
   d := &DSU{
       parent: make([]int, n),
       rank:   make([]int, n),
   }
   for i := 0; i < n; i++ {
       d.parent[i] = i
   }
   return d
}

// Find returns the representative of x
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union merges sets containing x and y
func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return
   }
   // union by rank
   if d.rank[rx] < d.rank[ry] {
       d.parent[rx] = ry
   } else if d.rank[ry] < d.rank[rx] {
       d.parent[ry] = rx
   } else {
       d.parent[ry] = rx
       d.rank[rx]++
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   colors := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &colors[i])
   }
   dsu := NewDSU(n)
   used := make([]bool, n)
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l--
       r--
       dsu.Union(l, r)
       used[l] = true
       used[r] = true
   }
   // For each component, count total nodes and color frequencies
   compCount := make(map[int]int)
   maxFreq := make(map[int]int)
   colorCount := make(map[int]map[int]int)
   for i := 0; i < n; i++ {
       if !used[i] {
           continue
       }
       root := dsu.Find(i)
       compCount[root]++
       col := colors[i]
       cmap, ok := colorCount[root]
       if !ok {
           cmap = make(map[int]int)
           colorCount[root] = cmap
       }
       cmap[col]++
       if cmap[col] > maxFreq[root] {
           maxFreq[root] = cmap[col]
       }
   }
   // Sum up minimal repaints
   result := 0
   for root, cnt := range compCount {
       // in each component, repaint all except the color with max frequency
       result += cnt - maxFreq[root]
   }
   fmt.Fprintln(writer, result)
}
