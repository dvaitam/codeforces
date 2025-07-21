package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1000000

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // find leaves (degree 1, not root)
   isLeaf := make([]bool, n+1)
   leafCount := 0
   for v := 2; v <= n; v++ {
       if len(adj[v]) == 1 {
           isLeaf[v] = true
           leafCount++
       }
   }
   // read leaf order
   leafOrder := make([]int, leafCount)
   for i := 0; i < leafCount; i++ {
       fmt.Fscan(reader, &leafOrder[i])
   }
   // map leaf to position
   pos := make([]int, n+1)
   for i, v := range leafOrder {
       pos[v] = i
   }
   // build rooted tree
   parent := make([]int, n+1)
   children := make([][]int, n+1)
   // BFS or DFS stack
   stack := []int{1}
   parent[1] = -1
   for len(stack) > 0 {
       v := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, u := range adj[v] {
           if u == parent[v] {
               continue
           }
           parent[u] = v
           children[v] = append(children[v], u)
           stack = append(stack, u)
       }
   }
   minpos := make([]int, n+1)
   maxpos := make([]int, n+1)
   cnt := make([]int, n+1)
   // init
   for i := 1; i <= n; i++ {
       minpos[i] = INF
       maxpos[i] = -INF
   }
   // post-order traversal list
   order := make([]int, 0, n)
   // simple stack for post-order
   var dfs func(int)
   dfs = func(v int) {
       for _, u := range children[v] {
           dfs(u)
       }
       order = append(order, v)
   }
   dfs(1)
   // compute
   for _, v := range order {
       if isLeaf[v] {
           cnt[v] = 1
           minpos[v] = pos[v]
           maxpos[v] = pos[v]
       }
       for _, u := range children[v] {
           cnt[v] += cnt[u]
           if minpos[u] < minpos[v] {
               minpos[v] = minpos[u]
           }
           if maxpos[u] > maxpos[v] {
               maxpos[v] = maxpos[u]
           }
       }
       if cnt[v] > 0 && maxpos[v]-minpos[v]+1 != cnt[v] {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // sort children by minpos
   for v := 1; v <= n; v++ {
       sort.Slice(children[v], func(i, j int) bool {
           return minpos[children[v][i]] < minpos[children[v][j]]
       })
   }
   // build answer
   res := make([]int, 0, 2*n-1)
   var build func(int)
   build = func(v int) {
       for _, u := range children[v] {
           res = append(res, u)
           build(u)
           res = append(res, v)
       }
   }
   res = append(res, 1)
   build(1)
   if len(res) != 2*n-1 {
       // should not happen
       fmt.Fprintln(writer, -1)
       return
   }
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
