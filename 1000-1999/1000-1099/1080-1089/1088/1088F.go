package main

import (
   "bufio"
   "fmt"
   "os"
)

const M = 21

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   w := make([]int, n+1)
   p := 1
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &w[i])
       if w[i] < w[p] {
           p = i
       }
   }
   adj := make([][]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(in, &u, &v)
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   // a[depth][i]: sparse table of mins on path
   a := make([][M]int, n)
   var ans int64
   // iterative DFS stack
   type frame struct{ u, parent, depth, next int }
   stack := make([]frame, 0, n)
   stack = append(stack, frame{u: p, parent: 0, depth: 0, next: -1})
   for len(stack) > 0 {
       f := &stack[len(stack)-1]
       if f.next == -1 {
           // first time visit
           dep := f.depth
           u := f.u
           a[dep][0] = w[u]
           for i := 1; i < M; i++ {
               a[dep][i] = a[dep][i-1]
               step := 1 << (i - 1)
               if dep >= step {
                   v := a[dep-step][i-1]
                   if v < a[dep][i] {
                       a[dep][i] = v
                   }
               }
           }
           if dep != 0 {
               ans += int64(w[u])
               nxt := a[dep-1][0]
               if dep > 1 {
                   now := dep - 2
                   for i := 0; i < M; i++ {
                       cost := a[now][i] * (i + 2)
                       if cost < nxt {
                           nxt = cost
                       }
                       now -= 1 << i
                       if now < 0 {
                           break
                       }
                   }
               }
               ans += int64(nxt)
           }
           f.next = 0
       }
       // process children
       u := f.u
       if f.next < len(adj[u]) {
           v := adj[u][f.next]
           f.next++
           if v == f.parent {
               continue
           }
           // push child
           stack = append(stack, frame{u: v, parent: u, depth: f.depth + 1, next: -1})
       } else {
           // done with this node
           stack = stack[:len(stack)-1]
       }
   }
   fmt.Fprintln(out, ans)
}
