package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   graph := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       graph[u] = append(graph[u], v)
       graph[v] = append(graph[v], u)
   }
   parent := make([]int, n+1)
   depth := make([]int64, n+1)
   order := make([]int, 0, n)
   // build parent, depth, and preorder
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   parent[1] = 0
   depth[1] = 0
   for i := 0; i < len(stack); i++ {
       x := stack[i]
       order = append(order, x)
       for _, y := range graph[x] {
           if y == parent[x] {
               continue
           }
           parent[y] = x
           depth[y] = depth[x] + 1
           stack = append(stack, y)
       }
   }
   val := make([]int64, n+1)
   // post-order accumulate subtree sums
   for i := len(order) - 1; i >= 0; i-- {
       x := order[i]
       val[x] = a[x]
       for _, y := range graph[x] {
           if y == parent[x] {
               continue
           }
           val[x] += val[y]
       }
   }
   dp := make([]int64, n+1)
   // initial dp at root 1
   for i := 1; i <= n; i++ {
       dp[1] += depth[i] * a[i]
   }
   // reroot dp
   for _, x := range order {
       for _, y := range graph[x] {
           if y == parent[x] {
               continue
           }
           dp[y] = dp[x] + val[1] - 2*val[y]
       }
   }
   // find max
   res := dp[1]
   for i := 2; i <= n; i++ {
       if dp[i] > res {
           res = dp[i]
       }
   }
   fmt.Println(res)
}
