package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU represents Disjoint Set Union (Union-Find) structure.
type DSU struct {
   p []int
}

// NewDSU initializes a DSU for n elements (0 to n-1).
func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   return &DSU{p: p}
}

// Find returns the representative of element x with path compression.
func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

// Union merges the sets containing a and b.
func (d *DSU) Union(a, b int) {
   ra := d.Find(a)
   rb := d.Find(b)
   if ra != rb {
       d.p[ra] = rb
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m1, m2 int
   fmt.Fscan(reader, &n, &m1, &m2)

   d1 := NewDSU(n)
   d2 := NewDSU(n)

   for i := 0; i < m1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       d1.Union(u, v)
   }
   for i := 0; i < m2; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       d2.Union(u, v)
   }

   var ans [][2]int
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if d1.Find(i) != d1.Find(j) && d2.Find(i) != d2.Find(j) {
               ans = append(ans, [2]int{i + 1, j + 1})
               d1.Union(i, j)
               d2.Union(i, j)
           }
       }
   }

   fmt.Fprintln(writer, len(ans))
   for _, p := range ans {
       fmt.Fprintln(writer, p[0], p[1])
   }
}
