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

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       // build graph
       adj := make([][]int, n+1)
       deg := make([]int, n+1)
       for i := 0; i < n; i++ {
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
           deg[u]++
           deg[v]++
       }
       // find cycle nodes by pruning leaves
       removed := make([]bool, n+1)
       queue := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if deg[i] == 1 {
               queue = append(queue, i)
           }
       }
       qi := 0
       for qi < len(queue) {
           u := queue[qi]
           qi++
           removed[u] = true
           for _, v := range adj[u] {
               if removed[v] {
                   continue
               }
               deg[v]--
               if deg[v] == 1 {
                   queue = append(queue, v)
               }
           }
       }
       iscycle := make([]bool, n+1)
       cycleRoots := make([]int, 0, n)
       for i := 1; i <= n; i++ {
           if !removed[i] {
               iscycle[i] = true
               cycleRoots = append(cycleRoots, i)
           }
       }
       // compute sizes s_i for each cycle root
       visited := make([]bool, n+1)
       var sumS2 int64
       for _, u := range cycleRoots {
           // DFS on tree hanging at u (excluding cycle edges)
           var cnt int64 = 1
           // stack for iterative DFS
           stack := make([]int, 0)
           for _, v := range adj[u] {
               if !iscycle[v] && !visited[v] {
                   visited[v] = true
                   stack = append(stack, v)
                   cnt++
               }
           }
           for len(stack) > 0 {
               v := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               for _, w := range adj[v] {
                   if !iscycle[w] && !visited[w] {
                       visited[w] = true
                       stack = append(stack, w)
                       cnt++
                   }
               }
           }
           sumS2 += cnt * cnt
       }
       // total paths = C(n,2) + sum_{i<j} s_i*s_j = n*(n-1)/2 + (n*n - sumS2)/2
       nn := int64(n)
       t1 := nn * (nn - 1) / 2
       t2 := (nn*nn - sumS2) / 2
       ans := t1 + t2
       fmt.Fprintln(writer, ans)
   }
}
