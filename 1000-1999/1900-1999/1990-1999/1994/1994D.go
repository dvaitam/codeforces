package main

import (
   "bufio"
   "fmt"
   "os"
)

// Disjoint set union (union-find)
type DSU struct {
   parent []int
}

// NewDSU initializes DSU for n elements
func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   return &DSU{parent: p}
}

// Find returns representative of x with path compression
func (d *DSU) Find(x int) int {
   if d.parent[x] != x {
       d.parent[x] = d.Find(d.parent[x])
   }
   return d.parent[x]
}

// Union merges sets containing a and b
func (d *DSU) Union(a, b int) {
   ra := d.Find(a)
   rb := d.Find(b)
   if ra != rb {
       d.parent[rb] = ra
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var T int
   if _, err := fmt.Fscan(in, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }

       dsu := NewDSU(n)
       // store edges (u, v)
       edges := make([][2]int, 0, n-1)
       // for t = n-1 down to 1
       for t := n - 1; t >= 1; t-- {
           mita := make([]int, t)
           for i := 0; i < t; i++ {
               mita[i] = -1
           }
           for i := 0; i < n; i++ {
               if dsu.Find(i) == i {
                   r := a[i] % t
                   if mita[r] == -1 {
                       mita[r] = i
                   } else {
                       u := mita[r]
                       v := i
                       edges = append(edges, [2]int{u, v})
                       dsu.Union(u, v)
                       break
                   }
               }
           }
       }
       // reverse edges to match order from t=1 to n-1
       for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
           edges[i], edges[j] = edges[j], edges[i]
       }

       // output
       fmt.Fprintln(out, "Yes")
       for _, e := range edges {
           // output 1-based indices
           fmt.Fprintln(out, e[0]+1, e[1]+1)
       }
   }
}
