package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   adj := make([][]int, n)
   parent := make([]int, n)
   for i := 2; i <= n; i++ {
       var y int
       fmt.Fscan(reader, &y)
       u, v := i-1, y-1
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // Build DFS order and parent
   order := make([]int, 0, n)
   stack := make([]int, 0, n)
   stack = append(stack, 0)
   parent[0] = -1
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
   // Compute subtree leaf counts
   s := make([]int, n)
   for i := n - 1; i >= 0; i-- {
       u := order[i]
       // count children
       cnt := len(adj[u])
       if parent[u] != -1 {
           cnt--
       }
       if cnt == 0 {
           s[u] = 1
       } else {
           sum := 0
           for _, v := range adj[u] {
               if v == parent[u] {
                   continue
               }
               sum += s[v]
           }
           s[u] = sum
       }
   }
   // sort and output
   sort.Ints(s)
   for i, v := range s {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
}
