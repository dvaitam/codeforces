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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   adj := make([][]int, n)
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   x := make([]bool, n)
   for i := 0; i < n; i++ {
       var xi int
       fmt.Fscan(reader, &xi)
       x[i] = (xi % 2) == 1
   }
   // find a root with x[root] == true
   root := -1
   for i := 0; i < n; i++ {
       if x[i] {
           root = i
           break
       }
   }
   if root == -1 {
       // all even, empty path
       fmt.Fprintln(writer, 0)
       return
   }
   visited := make([]bool, n)
   cur := make([]bool, n)
   ans := make([]int, 0, 4*n)
   // iterative DFS
   type frame struct{ u, parent, idx int }
   stack := []frame{{root, -1, 0}}
   for len(stack) > 0 {
       fr := &stack[len(stack)-1]
       u := fr.u
       if fr.idx == 0 {
           visited[u] = true
           ans = append(ans, u)
           cur[u] = !cur[u]
       }
       if fr.idx < len(adj[u]) {
           v := adj[u][fr.idx]
           fr.idx++
           if !visited[v] {
               stack = append(stack, frame{v, u, 0})
           }
           continue
       }
       // processed all children
       if cur[u] != x[u] {
           if fr.parent != -1 {
               // flip parent and back
               p := fr.parent
               ans = append(ans, p)
               cur[p] = !cur[p]
               ans = append(ans, u)
               cur[u] = !cur[u]
           } else {
               // root
               if len(ans) > 0 {
                   first := ans[0]
                   ans = ans[1:]
                   cur[first] = !cur[first]
               }
           }
       }
       // pop
       stack = stack[:len(stack)-1]
   }
   // check unreachable odd nodes
   for i := 0; i < n; i++ {
       if x[i] && !visited[i] {
           fmt.Fprintln(writer, -1)
           return
       }
   }
   // output
   k := len(ans)
   fmt.Fprintln(writer, k)
   if k > 0 {
       for i, v := range ans {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v+1)
       }
       fmt.Fprintln(writer)
   }
}
