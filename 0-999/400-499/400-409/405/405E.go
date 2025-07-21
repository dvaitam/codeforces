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
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   if m%2 == 1 {
       fmt.Fprintln(out, "No solution")
       return
   }
   u := make([]int, m)
   v := make([]int, m)
   adj := make([][]struct{to, eid int}, n+1)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &u[i], &v[i])
       adj[u[i]] = append(adj[u[i]], struct{to, eid int}{v[i], i})
       adj[v[i]] = append(adj[v[i]], struct{to, eid int}{u[i], i})
   }
   visited := make([]bool, n+1)
   used := make([]bool, m)
   indegree := make([]int, n+1)
   from := make([]int, m)
   to := make([]int, m)
   incoming := make([][]int, n+1)
   type stNode struct{u, parent, peid, idx int}
   stack := make([]stNode, 0, n)
   // start DFS at 1
   visited[1] = true
   stack = append(stack, stNode{u: 1, parent: -1, peid: -1, idx: 0})
   for len(stack) > 0 {
       node := &stack[len(stack)-1]
       u0 := node.u
       if node.idx < len(adj[u0]) {
           e := adj[u0][node.idx]
           node.idx++
           if used[e.eid] {
               continue
           }
           used[e.eid] = true
           v0 := e.to
           if !visited[v0] {
               visited[v0] = true
               stack = append(stack, stNode{u: v0, parent: u0, peid: e.eid, idx: 0})
           } else {
               // back-edge or cross-edge
               from[e.eid] = u0
               to[e.eid] = v0
               indegree[v0]++
               incoming[v0] = append(incoming[v0], e.eid)
           }
       } else {
           // backtrack
           if node.parent != -1 {
               v0 := node.u
               upar := node.parent
               eid := node.peid
               if indegree[v0]&1 == 1 {
                   // orient upar -> v0
                   from[eid] = upar
                   to[eid] = v0
                   indegree[v0]++
                   incoming[v0] = append(incoming[v0], eid)
               } else {
                   // orient v0 -> upar
                   from[eid] = v0
                   to[eid] = upar
                   indegree[upar]++
                   incoming[upar] = append(incoming[upar], eid)
               }
           }
           stack = stack[:len(stack)-1]
       }
   }
   // build and print paths
   w := bufio.NewWriter(out)
   defer w.Flush()
   for mid := 1; mid <= n; mid++ {
       lst := incoming[mid]
       for i := 0; i+1 < len(lst); i += 2 {
           e1 := lst[i]
           e2 := lst[i+1]
           x := from[e1]
           y := from[e2]
           fmt.Fprintln(w, x, mid, y)
       }
   }
