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
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // adjacency list
   g := make([][]int, n+1)
   lc := make([]int, n+1)
   lr := make([]int, n+1)

   for i := 1; i <= n; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       if lc[x] != 0 {
           u := lc[x]
           g[u] = append(g[u], i)
           g[i] = append(g[i], u)
           lc[x] = 0
       } else {
           lc[x] = i
       }
       if lr[y] != 0 {
           u := lr[y]
           g[u] = append(g[u], i)
           g[i] = append(g[i], u)
           lr[y] = 0
       } else {
           lr[y] = i
       }
   }

   col := make([]int, n+1)
   for i := 1; i <= n; i++ {
       col[i] = -1
   }
   // iterative DFS to color components
   stack := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if col[i] != -1 {
           continue
       }
       col[i] = 0
       stack = append(stack, i)
       for len(stack) > 0 {
           u := stack[len(stack)-1]
           stack = stack[:len(stack)-1]
           cu := col[u]
           for _, v := range g[u] {
               if col[v] == -1 {
                   col[v] = cu ^ 1
                   stack = append(stack, v)
               }
           }
       }
   }
   // output colors: 0 -> b, 1 -> r
   res := make([]byte, n)
   for i := 1; i <= n; i++ {
       if col[i] == 1 {
           res[i-1] = 'r'
       } else {
           res[i-1] = 'b'
       }
   }
   fmt.Fprint(out, string(res))
}
