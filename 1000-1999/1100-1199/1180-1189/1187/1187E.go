package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   parent[1] = 0
   // post-order using two stacks
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   order := make([]int, 0, n)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       order = append(order, u)
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           parent[v] = u
           stack = append(stack, v)
       }
   }
   // compute subtree sizes and dp for root 1
   size := make([]int, n+1)
   dp1 := int64(0)
   // order currently preorder, reverse for postorder
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       size[u] = 1
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           size[u] += size[v]
       }
       dp1 += int64(size[u])
   }
   // reroot dp
   dp := make([]int64, n+1)
   dp[1] = dp1
   ans := dp1
   // DFS to propagate
   stack = stack[:0]
   stack = append(stack, 1)
   for len(stack) > 0 {
       u := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           // move root from u to v
           dp[v] = dp[u] - int64(size[v]) + int64(n - size[v])
           if dp[v] > ans {
               ans = dp[v]
           }
           stack = append(stack, v)
       }
   }
   fmt.Fprintln(out, ans)
}
