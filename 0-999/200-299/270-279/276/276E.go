package main

import (
   "bufio"
   "fmt"
   "os"
)

// BIT supports range add and point query (difference array BIT)
type BIT struct {
   n int
   t []int64
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, t: make([]int64, n+1)}
}

// add v at position i (1-based)
func (b *BIT) add(i int, v int64) {
   for ; i <= b.n; i += i & -i {
       b.t[i] += v
   }
}

// sum returns prefix sum up to i (1-based)
func (b *BIT) sum(i int) int64 {
   var s int64
   for ; i > 0; i -= i & -i {
       s += b.t[i]
   }
   return s
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, q int
   fmt.Fscan(in, &n, &q)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // BFS from 1 to get parent, depth, branch assignment
   dep := make([]int, n+1)
   par := make([]int, n+1)
   branch := make([]int, n+1)
   dep[1] = 0
   branch[1] = 0
   // queue
   qn := make([]int, 0, n)
   qn = append(qn, 1)
   par[1] = 0
   // track max depth per branch
   maxdep := make([]int, n+1)
   for i := 0; i < len(qn); i++ {
       u := qn[i]
       for _, v := range adj[u] {
           if v == par[u] {
               continue
           }
           par[v] = u
           dep[v] = dep[u] + 1
           if u == 1 {
               branch[v] = v
           } else {
               branch[v] = branch[u]
           }
           if branch[v] != 0 && dep[v] > maxdep[branch[v]] {
               maxdep[branch[v]] = dep[v]
           }
           qn = append(qn, v)
       }
   }
   // prepare branch base offsets
   base := make([]int, n+1)
   total := 0
   // branches are children of 1
   for _, b := range adj[1] {
       // b is branch id
       sz := maxdep[b]
       if sz <= 0 {
           continue
       }
       base[b] = total + 1
       total += sz
   }
   // BIT for branch-specific, size = total
   bitBranch := NewBIT(total + 2)
   // BIT for global depth, max depth overall is max dep in dep[]
   maxd := 0
   for i := 1; i <= n; i++ {
       if dep[i] > maxd {
           maxd = dep[i]
       }
   }
   bitDepth := NewBIT(maxd + 3)

   // helper to get position in bitBranch
   pos := make([]int, n+1)
   for u := 2; u <= n; u++ {
       b := branch[u]
       if b != 0 {
           // depth starts at 1 for branch b
           pos[u] = base[b] + dep[u] - 1
       }
   }

   // process queries
   for i := 0; i < q; i++ {
       var typ int
       fmt.Fscan(in, &typ)
       if typ == 0 {
           var v, x, d0 int
           fmt.Fscan(in, &v, &x, &d0)
           dv := dep[v]
           // global updates for branch != v
           D2 := d0 - dv
           if D2 >= 0 {
               // depths 0..D2: map to indices 1..D2+1
               bitDepth.add(1, int64(x))
               bitDepth.add(D2+2, int64(-x))
           }
           // branch-specific for same branch
           if v != 1 {
               b := branch[v]
               // L0 = max(dep[v]-d0, 1), R0 = min(dep[v]+d0, maxdep[b])
               lo := dv - d0
               if lo < 1 {
                   lo = 1
               }
               hi := dv + d0
               if hi > maxdep[b] {
                   hi = maxdep[b]
               }
               // exclude depths <= D2 (global applied there)
               lb := lo
               if D2+1 > lb {
                   lb = D2 + 1
               }
               if lb <= hi {
                   lpos := base[b] + lb - 1
                   rpos := base[b] + hi - 1
                   bitBranch.add(lpos, int64(x))
                   bitBranch.add(rpos+1, int64(-x))
               }
           }
       } else {
           var v int
           fmt.Fscan(in, &v)
           // query v
           res := bitDepth.sum(dep[v] + 1)
           if branch[v] != 0 {
               res += bitBranch.sum(pos[v])
           }
           fmt.Fprintln(out, res)
       }
   }
}
