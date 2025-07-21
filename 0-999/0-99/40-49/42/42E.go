package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type DSU struct {
   p, r []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
   }
   return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) bool {
   rx, ry := d.Find(x), d.Find(y)
   if rx == ry {
       return false
   }
   if d.r[rx] < d.r[ry] {
       d.p[rx] = ry
   } else if d.r[ry] < d.r[rx] {
       d.p[ry] = rx
   } else {
       d.p[ry] = rx
       d.r[rx]++
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   fmt.Fscan(reader, &m)
   edges := make([][3]int, m)
   for i := 0; i < m; i++ {
       var u, v, c int
       fmt.Fscan(reader, &u, &v, &c)
       edges[i][0] = c
       edges[i][1] = u - 1
       edges[i][2] = v - 1
   }
   sort.Slice(edges, func(i, j int) bool {
       return edges[i][0] < edges[j][0]
   })
   dsu := NewDSU(n)
   var total int64
   cnt := 0
   for _, e := range edges {
       if dsu.Union(e[1], e[2]) {
           total += int64(e[0])
           cnt++
           if cnt == n-1 {
               break
           }
       }
   }
   // If not fully connected, no solution for any query
   var ans int64
   if cnt == n-1 {
       ans = total
   } else {
       ans = -1
   }
   var q int
   fmt.Fscan(reader, &q)
   for i := 0; i < q; i++ {
       // read ai, bi but unused
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if ans < 0 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, ans)
       }
   }
}
