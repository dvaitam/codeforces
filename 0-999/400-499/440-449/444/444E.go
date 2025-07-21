package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU represents a disjoint set union structure
type DSU struct {
   p []int
}

// NewDSU creates a DSU of size n
func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = -1
   }
   return &DSU{p: p}
}

// Find returns the root of x
func (d *DSU) Find(x int) int {
   if d.p[x] < 0 {
       return x
   }
   d.p[x] = d.Find(d.p[x])
   return d.p[x]
}

// Union merges sets of x and y
func (d *DSU) Union(x, y int) {
   x = d.Find(x)
   y = d.Find(y)
   if x == y {
       return
   }
   // union by size
   if d.p[x] > d.p[y] {
       x, y = y, x
   }
   d.p[x] += d.p[y]
   d.p[y] = x
}

// Edge represents a tree edge
type Edge struct {
   u, v, w int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   edges := make([]Edge, n-1)
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       edges[i] = Edge{u: a - 1, v: b - 1, w: c}
   }
   x := make([]int, n)
   totalCap := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &x[i])
       totalCap += x[i]
   }
   // collect weights
   W := make([]int, 0, n+1)
   W = append(W, 0)
   for _, e := range edges {
       W = append(W, e.w)
   }
   sort.Ints(W)
   // unique
   m := 0
   for i := 0; i < len(W); i++ {
       if i == 0 || W[i] != W[i-1] {
           W[m] = W[i]
           m++
       }
   }
   W = W[:m]
   // sort edges by weight
   sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
   // binary search on W
   lo, hi := 0, len(W)-1
   for lo < hi {
       mid := (lo + hi + 1) / 2
       if check(W[mid], n, edges, x, totalCap) {
           lo = mid
       } else {
           hi = mid - 1
       }
   }
   // output answer
   fmt.Println(W[lo])
}

// check returns true if threshold T is feasible
func check(T, n int, edges []Edge, x []int, totalCap int) bool {
   dsu := NewDSU(n)
   // union edges with weight < T
   for _, e := range edges {
       if e.w >= T {
           break
       }
       dsu.Union(e.u, e.v)
   }
   // accumulate component sizes and capacities
   compSize := make(map[int]int)
   compCap := make(map[int]int)
   for i := 0; i < n; i++ {
       r := dsu.Find(i)
       compSize[r]++
       compCap[r] += x[i]
   }
   // check Hall's condition per component
   for r, sz := range compSize {
       capc := compCap[r]
       if sz+capc > totalCap {
           return false
       }
   }
   return true
}
