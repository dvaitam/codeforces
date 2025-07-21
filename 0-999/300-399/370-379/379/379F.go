package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   // max nodes = initial 4 + 2*q
   maxn := 2*q + 5
   // number of levels for binary lifting
   L := bits.Len(uint(maxn))

   // parent[k][v]: the 2^k-th ancestor of v
   parent := make([][]int32, L)
   for i := 0; i < L; i++ {
       parent[i] = make([]int32, maxn)
   }
   depth := make([]int, maxn)
   // initial tree: 1 connected to 2,3,4
   // root 1
   depth[1] = 0
   parent[0][1] = 0
   for v := 2; v <= 4; v++ {
       depth[v] = 1
       parent[0][v] = 1
   }
   // fill ancestors for initial nodes
   for k := 1; k < L; k++ {
       for v := 1; v <= 4; v++ {
           p := parent[k-1][v]
           parent[k][v] = parent[k-1][p]
       }
   }
   // current endpoints of diameter
   a, b := 2, 3
   curD := 2
   n := 4

   // lca function
   lca := func(u, v int) int {
       if depth[u] < depth[v] {
           u, v = v, u
       }
       // lift u to depth v
       diff := depth[u] - depth[v]
       for k := 0; diff > 0; k++ {
           if diff&1 != 0 {
               u = int(parent[k][u])
           }
           diff >>= 1
       }
       if u == v {
           return u
       }
       for k := L - 1; k >= 0; k-- {
           pu := parent[k][u]
           pv := parent[k][v]
           if pu != pv {
               u = int(pu)
               v = int(pv)
           }
       }
       return int(parent[0][u])
   }
   // distance
   dist := func(u, v int) int {
       w := lca(u, v)
       return depth[u] + depth[v] - 2*depth[w]
   }

   // process queries
   for i := 0; i < q; i++ {
       var v int
       fmt.Fscan(in, &v)
       // add two new leaves u and w
       // first
       n++
       u := n
       depth[u] = depth[v] + 1
       parent[0][u] = int32(v)
       for k := 1; k < L; k++ {
           p := parent[k-1][u]
           parent[k][u] = parent[k-1][p]
       }
       // update diameter with u
       du := dist(a, u)
       if du > curD {
           curD = du
           b = u
       } else {
           dv := dist(b, u)
           if dv > curD {
               curD = dv
               a = u
           }
       }
       // second
       n++
       w := n
       depth[w] = depth[v] + 1
       parent[0][w] = int32(v)
       for k := 1; k < L; k++ {
           p := parent[k-1][w]
           parent[k][w] = parent[k-1][p]
       }
       du = dist(a, w)
       if du > curD {
           curD = du
           b = w
       } else {
           dv := dist(b, w)
           if dv > curD {
               curD = dv
               a = w
           }
       }
       // output
       fmt.Fprintln(out, curD)
   }
}
