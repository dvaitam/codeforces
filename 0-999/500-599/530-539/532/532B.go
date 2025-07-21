package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // tree
   children := make([][]int, n+1)
   a := make([]int64, n+1)
   var p int
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &p, &a[i])
       if p >= 1 {
           children[p] = append(children[p], i)
       }
   }
   // dp0[u]: max sum with even selected in subtree u
   // dp1[u]: max sum with odd selected in subtree u
   dp0 := make([]int64, n+1)
   dp1 := make([]int64, n+1)
   const inf64 = int64(4e18)
   negInf := -inf64

   // post-order traversal
   type st struct{ u, state int }
   stack := make([]st, 0, n*2)
   stack = append(stack, st{1, 0})
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       u, state := cur.u, cur.state
       if state == 0 {
           // enter
           stack = append(stack, st{u, 1})
           for _, v := range children[u] {
               stack = append(stack, st{v, 0})
           }
       } else {
           // exit: compute dp
           c0, c1 := int64(0), negInf
           for _, v := range children[u] {
               // merge child v
               nc0 := max(c0+dp0[v], c1+dp1[v])
               nc1 := max(c0+dp1[v], c1+dp0[v])
               c0, c1 = nc0, nc1
           }
           // u not selected: c0, c1
           // u selected: descendants must be even => use c0, adds a[u]
           sel1 := c0 + a[u]
           dp0[u] = c0
           dp1[u] = max(c1, sel1)
       }
   }
   // answer
   res := max(dp0[1], dp1[1])
   fmt.Fprintln(writer, res)
}
