package main

import (
   "bufio"
   "fmt"
   "os"
)

// Disjoint set union
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   d := &DSU{p: make([]int, n+1)}
   for i := 1; i <= n; i++ {
       d.p[i] = i
   }
   return d
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(a, b int) {
   pa := d.Find(a)
   pb := d.Find(b)
   if pa != pb {
       d.p[pa] = pb
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m1, m2 int
   if _, err := fmt.Fscan(reader, &n, &m1, &m2); err != nil {
       return
   }
   d1 := NewDSU(n)
   d2 := NewDSU(n)
   for i := 0; i < m1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       d1.Union(u, v)
   }
   for i := 0; i < m2; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       d2.Union(u, v)
   }
   type pair struct{ x, y int }
   var res []pair

   // connect to node 1 greedily
   for i := 2; i <= n; i++ {
       if d1.Find(i) != d1.Find(1) && d2.Find(i) != d2.Find(1) {
           res = append(res, pair{i, 1})
           d1.Union(i, 1)
           d2.Union(i, 1)
       }
   }
   // collect remaining component roots (excluding 1)
   var s1, s2 []int
   root1 := d1.Find(1)
   root2 := d2.Find(1)
   for i := 2; i <= n; i++ {
       if d1.Find(i) == i && d1.Find(i) != root1 {
           s1 = append(s1, i)
       }
       if d2.Find(i) == i && d2.Find(i) != root2 {
           s2 = append(s2, i)
       }
   }
   // match remaining
   t := len(s1)
   if len(s2) < t {
       t = len(s2)
   }
   for i := 0; i < t; i++ {
       res = append(res, pair{s1[i], s2[i]})
   }
   // output
   fmt.Fprintln(writer, len(res))
   for _, p := range res {
       fmt.Fprintln(writer, p.x, p.y)
   }
}
