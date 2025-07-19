package main

import (
   "fmt"
   "os"
)

func main() {
   var n, k int
   if _, err := fmt.Fscan(os.Stdin, &n, &k); err != nil {
       return
   }
   graph := make([][]int, n+1)
   deg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(os.Stdin, &u, &v)
       graph[u] = append(graph[u], v)
       graph[v] = append(graph[v], u)
       deg[u]++
       deg[v]++
   }
   vis := make([]bool, n+1)
   d := make([]int, n+1)
   queue := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       if deg[i] == 1 {
           queue = append(queue, i)
           vis[i] = true
       }
   }
   head := 0
   ok := true
   for head < len(queue) && ok {
       x := queue[head]
       head++
       for _, v := range graph[x] {
           if vis[v] {
               if d[v] == d[x]-1 || d[v] == d[x]+1 {
                   continue
               }
               ok = false
               break
           }
           d[v] = d[x] + 1
           queue = append(queue, v)
           vis[v] = true
       }
   }
   if !ok {
       fmt.Println("NO")
       return
   }
   center := queue[len(queue)-1]
   if d[center] != k || deg[center] < 3 {
       fmt.Println("NO")
       return
   }
   for i := 1; i <= n; i++ {
       if i == center {
           continue
       }
       if deg[i] != 1 && deg[i] < 4 {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
