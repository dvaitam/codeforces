package main

import (
   "bufio"
   "fmt"
   "os"
)

// DSU (Disjoint Set Union) for integers 0..n-1
type DSU struct {
   p, r []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   r := make([]int, n)
   for i := 0; i < n; i++ {
       p[i] = i
       r[i] = 0
   }
   return &DSU{p: p, r: r}
}

func (d *DSU) Find(x int) int {
   if d.p[x] != x {
       d.p[x] = d.Find(d.p[x])
   }
   return d.p[x]
}

func (d *DSU) Union(x, y int) {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return
   }
   if d.r[rx] < d.r[ry] {
       d.p[rx] = ry
   } else if d.r[ry] < d.r[rx] {
       d.p[ry] = rx
   } else {
       d.p[ry] = rx
       d.r[rx]++
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   p := make([]int, n)
   q := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &p[i])
       p[i]--
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &q[i])
       q[i]--
   }
   dsu := NewDSU(n)
   // union adjacent in p and q
   for i := 0; i+1 < n; i++ {
       dsu.Union(p[i], p[i+1])
       dsu.Union(q[i], q[i+1])
   }
   // assign component IDs
   compID := make([]int, n)
   for i := range compID {
       compID[i] = -1
   }
   comps := 0
   for i := 0; i < n; i++ {
       r := dsu.Find(i)
       if compID[r] == -1 {
           compID[r] = comps
           comps++
       }
       compID[i] = compID[r]
   }
   if comps < k {
       fmt.Fprintln(out, "NO")
       return
   }
   // build graph of components
   indeg := make([]int, comps)
   adj := make([][]int, comps)
   for i := 0; i+1 < n; i++ {
       u := compID[p[i]]
       v := compID[p[i+1]]
       if u != v {
           adj[u] = append(adj[u], v)
           indeg[v]++
       }
       u = compID[q[i]]
       v = compID[q[i+1]]
       if u != v {
           adj[u] = append(adj[u], v)
           indeg[v]++
       }
   }
   // Kahn's algorithm for topological sort
   order := make([]int, 0, comps)
   queue := make([]int, 0, comps)
   for i := 0; i < comps; i++ {
       if indeg[i] == 0 {
           queue = append(queue, i)
       }
   }
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       order = append(order, u)
       for _, v := range adj[u] {
           indeg[v]--
           if indeg[v] == 0 {
               queue = append(queue, v)
           }
       }
   }
   // assign letters to components according to topo order
   letters := make([]byte, comps)
   for i, comp := range order {
       if i < k {
           letters[comp] = byte('a' + i)
       } else {
           letters[comp] = byte('a' + k - 1)
       }
   }
   // build result string
   res := make([]byte, n)
   for i := 0; i < n; i++ {
       res[i] = letters[compID[i]]
   }
   fmt.Fprintln(out, "YES")
   out.Write(res)
   out.WriteByte('\n')
}
