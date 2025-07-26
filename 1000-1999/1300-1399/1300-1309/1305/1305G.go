package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU with path compression
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   p := make([]int, n)
   for i := range p {
       p[i] = -1
   }
   return &DSU{p}
}

func (d *DSU) Find(x int) int {
   if d.p[x] < 0 {
       return x
   }
   d.p[x] = d.Find(d.p[x])
   return d.p[x]
}

func (d *DSU) Union(a, b int) bool {
   a = d.Find(a)
   b = d.Find(b)
   if a == b {
       return false
   }
   if d.p[a] > d.p[b] {
       a, b = b, a
   }
   d.p[a] += d.p[b]
   d.p[b] = a
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // Prepare nodes sorted by age descending
   idx := make([]int, n)
   for i := range idx {
       idx[i] = i
   }
   sort.Slice(idx, func(i, j int) bool {
       return a[idx[i]] > a[idx[j]]
   })
   const bits = 18
   const fullMask = (1<<bits - 1)
   masks := make([]int, n)
   for i := 0; i < n; i++ {
       masks[i] = a[i]
   }
   dsu := NewDSU(n)
   active := make([]bool, n)
   size := 1 << bits
   ids := make([][]int, size)
   total := int64(0)
   // first node
   u0 := idx[0]
   active[u0] = true
   ids[masks[u0]] = append(ids[masks[u0]], u0)
   // process others
   for i := 1; i < n; i++ {
       u := idx[i]
       mu := masks[u]
       active[u] = true
       if mu == 0 {
           // connect to first node
           if dsu.Union(u, u0) {
               total += int64(a[u0])
           }
           ids[0] = append(ids[0], u)
           continue
       }
       M := fullMask ^ mu
       for s := M; ; s = (s - 1) & M {
           for _, v := range ids[s] {
               if active[v] && dsu.Union(u, v) {
                   // v was processed earlier, so a[v] >= a[u]
                   total += int64(a[v])
               }
           }
           if s == 0 {
               break
           }
       }
       ids[mu] = append(ids[mu], u)
   }
   // output result
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, total)
}
