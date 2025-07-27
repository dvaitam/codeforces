package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
      return
   }

   g := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
      var u, v int
      fmt.Fscan(reader, &u, &v)
      g[u] = append(g[u], v)
      g[v] = append(g[v], u)
   }

   parent := make([]int, n+1)
   depth := make([]int, n+1)
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   timer := 0

   type frame struct{ v, p, idx int }
   stack := []frame{{1, 0, -1}}
   for len(stack) > 0 {
      f := &stack[len(stack)-1]
      if f.idx == -1 {
         parent[f.v] = f.p
         depth[f.v] = depth[f.p] + 1
         timer++
         tin[f.v] = timer
         f.idx = 0
      }
      if f.idx < len(g[f.v]) {
         to := g[f.v][f.idx]
         f.idx++
         if to == f.p {
            continue
         }
         stack = append(stack, frame{to, f.v, -1})
      } else {
         timer++
         tout[f.v] = timer
         stack = stack[:len(stack)-1]
      }
   }

   isAncestor := func(u, v int) bool {
      return tin[u] <= tin[v] && tout[v] <= tout[u]
   }

   for ; m > 0; m-- {
      var k int
      fmt.Fscan(reader, &k)
      nodes := make([]int, k)
      deepest := 1
      for i := 0; i < k; i++ {
         fmt.Fscan(reader, &nodes[i])
         if nodes[i] != 1 {
            nodes[i] = parent[nodes[i]]
         }
         if depth[nodes[i]] > depth[deepest] {
            deepest = nodes[i]
         }
      }
      ok := true
      for _, v := range nodes {
         if !isAncestor(v, deepest) {
            ok = false
            break
         }
      }
      if ok {
         fmt.Fprintln(writer, "YES")
      } else {
         fmt.Fprintln(writer, "NO")
      }
   }
}
