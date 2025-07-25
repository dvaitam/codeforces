package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU for zero-edge components
type DSU struct {
   p, sz []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n+1)
   sz := make([]int, n+1)
   for i := 1; i <= n; i++ {
       p[i] = i
       sz[i] = 1
   }
   return &DSU{p: p, sz: sz}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   x = d.Find(x)
   y = d.Find(y)
   if x == y {
       return
   }
   // union by size
   if d.sz[x] < d.sz[y] {
       x, y = y, x
   }
   d.p[y] = x
   d.sz[x] += d.sz[y]
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   dsu := NewDSU(n)
   // count of nodes is n
   // process edges
   edges := make([][3]int, n-1)
   for i := 0; i < n-1; i++ {
       var u, v, c int
       fmt.Fscan(in, &u, &v, &c)
       edges[i] = [3]int{u, v, c}
       if c == 0 {
           dsu.Union(u, v)
       }
   }
   // compute component sizes and count components
   compSz := make(map[int]int)
   for i := 1; i <= n; i++ {
       r := dsu.Find(i)
       compSz[r]++
   }
   k := len(compSz)
   // sum within-component ordered pairs
   var ans int64
   for _, s := range compSz {
       ss := int64(s)
       ans += ss * (ss - 1)
   }
   // add cross-component valid pairs
   // for each x, there is one y in each other component: n * (k-1)
   ans += int64(n) * int64(k-1)
   // output
   fmt.Println(ans)
}
