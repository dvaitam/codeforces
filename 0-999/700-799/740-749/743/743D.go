package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n < 2 {
       fmt.Println("Impossible")
       return
   }
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   parent := make([]int, n+1)
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 1)
   parent[1] = 0
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
   sum := make([]int64, n+1)
   f := make([]int64, n+1)
   const INF = int64(4e18)
   var ans int64
   hasAns := false
   // process in reverse order for post-order
   for i := n - 1; i >= 0; i-- {
       u := order[i]
       sum[u] = a[u]
       // track top two best f[v]
       best1, best2 := -INF, -INF
       for _, v := range adj[u] {
           if v == parent[u] {
               continue
           }
           sum[u] += sum[v]
           // update best1, best2
           if f[v] > best1 {
               best2 = best1
               best1 = f[v]
           } else if f[v] > best2 {
               best2 = f[v]
           }
       }
       // if two disjoint subtrees here
       if best2 > -INF {
           cand := best1 + best2
           if !hasAns || cand > ans {
               ans = cand
               hasAns = true
           }
       }
       // f[u] is best subtree sum in this subtree
       f[u] = sum[u]
       if best1 > f[u] {
           f[u] = best1
       }
   }
   if !hasAns {
       fmt.Println("Impossible")
   } else {
       fmt.Println(ans)
   }
}
