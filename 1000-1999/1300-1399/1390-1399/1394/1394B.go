package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   // read edges and build outgoing lists
   out := make([][]pair, n+1)
   for e := 0; e < m; e++ {
       var u, v, w int
       fmt.Fscan(in, &u, &v, &w)
       out[u] = append(out[u], pair{w, v})
   }
   // sort outgoing by weight
   degCount := make([]int, k+1)
   // record destinations per category
   // first, count deg per node and sort edges
   for u := 1; u <= n; u++ {
       lst := out[u]
       sort.Slice(lst, func(i, j int) bool { return lst[i].w < lst[j].w })
       out[u] = lst
       d := len(lst)
       if d <= k {
           degCount[d]++
       }
   }
   // category id assignment
   // max categories = k*(k+1)/2
   maxCats := k*(k+1)/2
   catIdx := make([][]int, k+1)
   for i := 1; i <= k; i++ {
       catIdx[i] = make([]int, i+1)
       for j := 1; j <= i; j++ {
           catIdx[i][j] = -1
       }
   }
   dests := make([][]int, 0, maxCats)
   // build dest lists and filter invalid
   cid := 0
   for i := 1; i <= k; i++ {
       // for each j from 1..i, collect destinations of nodes with deg i
       for j := 1; j <= i; j++ {
           // if no node has deg i, this category still valid with empty
           var dlist []int
           if degCount[i] > 0 {
               dlist = make([]int, 0, degCount[i])
               // collect
               for u := 1; u <= n; u++ {
                   if len(out[u]) == i {
                       dlist = append(dlist, out[u][j-1].v)
                   }
               }
               sort.Ints(dlist)
               ok := true
               for p := 1; p < len(dlist); p++ {
                   if dlist[p] == dlist[p-1] {
                       ok = false
                       break
                   }
               }
               if !ok {
                   continue
               }
           }
           // assign category
           catIdx[i][j] = cid
           dests = append(dests, dlist)
           cid++
       }
   }
   C := cid
   // build conflict bitmask
   conflicts := make([]uint64, C)
   for a := 0; a < C; a++ {
       conflicts[a] = 1 << a
   }
   // map dest to categories
   // build conflict bitmask using destination lists
   destMap := make([][]int, n+1)
   for a := 0; a < C; a++ {
       for _, v := range dests[a] {
           destMap[v] = append(destMap[v], a)
       }
   }
   // for each destination, mark conflicts among categories sharing it
   for v := 1; v <= n; v++ {
       cats := destMap[v]
       if len(cats) < 2 {
           continue
       }
       for _, a := range cats {
           for _, b := range cats {
               if a != b {
                   conflicts[a] |= 1 << b
               }
           }
       }
   }
   // prepare types
   typeCats := make([][]int, k+1)
   for i := 1; i <= k; i++ {
       for j := 1; j <= i; j++ {
           if id := catIdx[i][j]; id >= 0 {
               typeCats[i] = append(typeCats[i], id)
           }
       }
   }
   // dfs over types
   var dfs func(int, uint64) int64
   dfs = func(i int, used uint64) int64 {
       if i > k {
           return 1
       }
       var cnt int64
       for _, id := range typeCats[i] {
           if used&conflicts[id] == 0 {
               cnt += dfs(i+1, used | (1<<id))
           }
       }
       return cnt
   }
   res := dfs(1, 0)
   fmt.Println(res)
}

type pair struct{ w, v int }
