package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}
func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, val int
   fmt.Fscan(reader, &n, &val)
   adj := make([][]int, n+1)
   deg := make([]int, n+1)
   for i := 2; i <= n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
       deg[u]++
       deg[v]++
   }
   // find root: first node with deg != 1
   root := 1
   for i := 1; i <= n; i++ {
       if deg[i] != 1 {
           root = i
           break
       }
   }
   parent := make([]int, n+1)
   // postorder traversal
   type item struct{ u, p int; visited bool }
   stack := make([]item, 0, n)
   stack = append(stack, item{root, 0, false})
   parent[root] = 0
   order := make([]int, 0, n)
   for len(stack) > 0 {
       it := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, p := it.u, it.p
       if !it.visited {
           // push node again for postorder
           stack = append(stack, item{u, p, true})
           // push children
           for _, v := range adj[u] {
               if v == p {
                   continue
               }
               parent[v] = u
               stack = append(stack, item{v, u, false})
           }
       } else {
           order = append(order, u)
       }
   }
   // DP f values
   const INF = 1e9
   f := make([]int, n+1)
   ans := 0
   for _, u := range order {
       if deg[u] == 1 {
           f[u] = 0
       } else {
           f[u] = -INF
       }
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           // combine f[u] and f[v]
           if f[u]+f[v]+1 > val {
               ans++
               f[u] = min(f[u], f[v]+1)
           } else {
               f[u] = max(f[u], f[v]+1)
           }
       }
   }
   // output partitions count
   fmt.Fprintln(writer, ans+1)
}
