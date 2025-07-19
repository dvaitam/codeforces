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
   sz := make([]int64, n+1)
   var cnt0, cnt1, ans int64

   type frame struct {
       u, parent, depth int
       visited          bool
   }
   stack := make([]frame, 0, 2*n)
   stack = append(stack, frame{u: 1, parent: 0, depth: 0, visited: false})
   for len(stack) > 0 {
       frm := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if !frm.visited {
           if frm.depth == 0 {
               cnt0++
           } else {
               cnt1++
           }
           stack = append(stack, frame{u: frm.u, parent: frm.parent, depth: frm.depth, visited: true})
           for _, v := range adj[frm.u] {
               if v != frm.parent {
                   stack = append(stack, frame{u: v, parent: frm.u, depth: frm.depth ^ 1, visited: false})
               }
           }
       } else {
           sum := int64(1)
           for _, v := range adj[frm.u] {
               if v != frm.parent {
                   sum += sz[v]
               }
           }
           sz[frm.u] = sum
           ans += sum * (int64(n) - sum)
       }
   }

   res := (ans + cnt0*cnt1) >> 1
   fmt.Fprintln(writer, res)
}
