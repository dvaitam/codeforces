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

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   totalSpecial := 2 * k
   special := make([]int, n+1)
   for i := 0; i < totalSpecial; i++ {
       var u int
       fmt.Fscan(reader, &u)
       special[u] = 1
   }
   adj := make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       adj[x] = append(adj[x], y)
       adj[y] = append(adj[y], x)
   }
   parent := make([]int, n+1)
   order := make([]int, 0, n)
   // iterative DFS to get parent and order
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
   cnt := make([]int, n+1)
   var result int64
   // process in post-order
   for i := len(order) - 1; i >= 0; i-- {
       u := order[i]
       cnt[u] += special[u]
       if u != 1 {
           x := cnt[u]
           // pairs crossing edge u-parent[u]
           result += int64(min(x, totalSpecial-x))
           cnt[parent[u]] += cnt[u]
       }
   }
   fmt.Fprintln(writer, result)
}
