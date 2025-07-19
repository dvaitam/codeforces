package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "sort"
)

func solve(g [][]int) []int {
   n := len(g)
   // random order
   order := make([]int, n)
   for i := 0; i < n; i++ {
       order[i] = i
   }
   for i := 1; i < n; i++ {
       j := rand.Intn(i + 1)
       order[i], order[j] = order[j], order[i]
   }
   take := make([]bool, n)
   reach := make([]bool, n)
   var taken []int
   for _, v := range order {
       if reach[v] {
           continue
       }
       valid := true
       for _, to := range g[v] {
           if take[to] {
               valid = false
               break
           }
       }
       if !valid {
           continue
       }
       take[v] = true
       taken = append(taken, v)
       for _, to := range g[v] {
           reach[to] = true
       }
   }
   term := true
   for v := 0; v < n; v++ {
       if !take[v] && !reach[v] {
           term = false
           break
       }
   }
   if term {
       return taken
   }
   // build subgraph of uncovered vertices
   newId := make([]int, n)
   for i := range newId {
       newId[i] = -1
   }
   var base []int
   cnt := 0
   for v := 0; v < n; v++ {
       if take[v] || reach[v] {
           continue
       }
       newId[v] = cnt
       base = append(base, v)
       cnt++
   }
   newG := make([][]int, cnt)
   for v := 0; v < n; v++ {
       idv := newId[v]
       if idv < 0 {
           continue
       }
       for _, to := range g[v] {
           if newId[to] >= 0 {
               newG[idv] = append(newG[idv], newId[to])
           }
       }
   }
   res := solve(newG)
   // map back
   for i := range res {
       res[i] = base[res[i]]
   }
   ok := make([]bool, n)
   for _, v := range res {
       for _, to := range g[v] {
           ok[to] = true
       }
   }
   for _, v := range taken {
       if !ok[v] {
           res = append(res, v)
       }
   }
   return res
}

func main() {
   rand.Seed(998244353)
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   g := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       g[u] = append(g[u], v)
   }
   ans := solve(g)
   sort.Ints(ans)
   fmt.Fprintln(writer, len(ans))
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v+1)
   }
   writer.WriteByte('\n')
}
