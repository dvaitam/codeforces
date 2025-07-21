package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Edge struct {
   u, v   int
   g, s   int64
   idxG   int
}

// DSU supports union-find
type DSU struct {
   p, r []int
}
func NewDSU(n int) *DSU {
   p := make([]int, n)
   r := make([]int, n)
   for i := range p {
       p[i] = i
   }
   return &DSU{p: p, r: r}
}
func (d *DSU) Find(x int) int {
   for d.p[x] != x {
       d.p[x] = d.p[d.p[x]]
       x = d.p[x]
   }
   return x
}
func (d *DSU) Union(x, y int) bool {
   rx := d.Find(x)
   ry := d.Find(y)
   if rx == ry {
       return false
   }
   if d.r[rx] < d.r[ry] {
       d.p[rx] = ry
   } else if d.r[rx] > d.r[ry] {
       d.p[ry] = rx
   } else {
       d.p[ry] = rx
       d.r[rx]++
   }
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M int
   var G, S int64
   if _, err := fmt.Fscan(in, &N, &M); err != nil {
       return
   }
   fmt.Fscan(in, &G, &S)
   edges := make([]Edge, M)
   for i := 0; i < M; i++ {
       var x, y int
       var gi, si int64
       fmt.Fscan(in, &x, &y, &gi, &si)
       edges[i] = Edge{u: x - 1, v: y - 1, g: gi, s: si}
   }
   // sort by g
   sort.Slice(edges, func(i, j int) bool {
       return edges[i].g < edges[j].g
   })
   for i := range edges {
       edges[i].idxG = i
   }
   // find minimal prefix i0 where edges[0..i0] connect graph by g
   d0 := NewDSU(N)
   comps := N
   i0 := -1
   for i, e := range edges {
       if d0.Union(e.u, e.v) {
           comps--
           if comps == 1 {
               i0 = i
               break
           }
       }
   }
   if i0 < 0 {
       fmt.Println(-1)
       return
   }
   // sort edges by s
   edgesS := make([]Edge, M)
   copy(edgesS, edges)
   sort.Slice(edgesS, func(i, j int) bool {
       return edgesS[i].s < edgesS[j].s
   })

   const INF = int64(4e18)
   ans := INF
   // for each prefix by g
   for i := i0; i < M; i++ {
       a := edges[i].g
       // prune by current ans
       if a*G >= ans {
           break
       }
       // find minimal B by si MST on edges with idxG <= i
       dsu := NewDSU(N)
       used := 0
       var b int64
       for _, e := range edgesS {
           if e.idxG > i {
               continue
           }
           if dsu.Union(e.u, e.v) {
               used++
               if e.s > b {
                   b = e.s
               }
               if used == N-1 {
                   break
               }
           }
       }
       if used == N-1 {
           cost := a*G + b*S
           if cost < ans {
               ans = cost
           }
       }
   }
   if ans == INF {
       fmt.Println(-1)
   } else {
       fmt.Println(ans)
   }
}
