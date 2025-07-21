package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Edge represents a pipe
type Edge struct{u, v, l, h, a int}
type E2 = Edge

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // read input edges
   orig := make([]Edge, 0, n*(n-1)/2)
   for i := 0; i < n*(n-1)/2; i++ {
       var u, v, l, h, a int
       fmt.Fscan(in, &u, &v, &l, &h, &a)
       orig = append(orig, Edge{u - 1, v - 1, l, h, a})
   }
   // sort edges by source
   edgelist := make([]*E2, len(orig))
   for i, e := range orig {
       edgelist[i] = &E2{e.u, e.v, e.l, e.h, e.a}
   }
   sort.Slice(edgelist, func(i, j int) bool {
       if edgelist[i].u != edgelist[j].u {
           return edgelist[i].u < edgelist[j].u
       }
       return edgelist[i].v < edgelist[j].v
   })
   src, sink := 0, n-1
   m := len(edgelist)
   // mark last edge for each source
   last := make([]bool, m)
   for i := 0; i < m; i++ {
       if i == m-1 || edgelist[i].u != edgelist[i+1].u {
           last[i] = true
       }
   }
   // compute flow bounds from source
   sumL, sumH := 0, 0
   for _, e := range edgelist {
       if e.u == src {
           sumL += e.l
           sumH += e.h
       }
   }
   // balances for nodes
   balance := make([]int, n)
   // dfs1: check feasibility for targetF
   var targetF int
   var dfs1 func(int) bool
   dfs1 = func(i int) bool {
       if i == m {
           F := -balance[src]
           return F == targetF && balance[sink] == F
       }
       e := edgelist[i]
       for c := e.l; c <= e.h; c++ {
           balance[e.u] -= c
           balance[e.v] += c
           prune := false
           if last[i] {
               u := e.u
               if u != src && u != sink && balance[u] != 0 {
                   prune = true
               }
               if u == src && -balance[src] != targetF {
                   prune = true
               }
           }
           if !prune && dfs1(i+1) {
               return true
           }
           balance[e.u] += c
           balance[e.v] -= c
       }
       return false
   }
   // find minimal feasible flow
   bestF := -1
   for f := sumL; f <= sumH; f++ {
       targetF = f
       for i := range balance {
           balance[i] = 0
       }
       if dfs1(0) {
           bestF = f
           break
       }
   }
   if bestF < 0 {
       fmt.Println("-1 -1")
       return
   }
   // dfs2: search maximal cost for bestF
   for i := range balance {
       balance[i] = 0
   }
   var bestCost int64 = -1
   var dfs2 func(int, int64)
   dfs2 = func(i int, cost int64) {
       if i == m {
           F := -balance[src]
           if F == bestF && balance[sink] == F && cost > bestCost {
               bestCost = cost
           }
           return
       }
       e := edgelist[i]
       for c := e.l; c <= e.h; c++ {
           balance[e.u] -= c
           balance[e.v] += c
           newCost := cost
           if c > 0 {
               newCost += int64(e.a) + int64(c)*int64(c)
           }
           prune := false
           if last[i] {
               u := e.u
               if u != src && u != sink && balance[u] != 0 {
                   prune = true
               }
               if u == src && -balance[src] != bestF {
                   prune = true
               }
           }
           if !prune {
               dfs2(i+1, newCost)
           }
           balance[e.u] += c
           balance[e.v] -= c
       }
   }
   dfs2(0, 0)
   fmt.Printf("%d %d\n", bestF, bestCost)
}
