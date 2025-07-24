package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   colors := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &colors[i])
   }
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       u--
       v--
       adj[u] = append(adj[u], v)
       adj[v] = append(adj[v], u)
   }
   comp := [2]int{}
   visited := make([]bool, n)
   stack := make([]int, 0, n)
   for target := 0; target <= 1; target++ {
       // reset visited for this color
       for i := 0; i < n; i++ {
           visited[i] = false
       }
       for i := 0; i < n; i++ {
           if !visited[i] && colors[i] == target {
               comp[target]++
               // DFS on same-color component
               stack = stack[:0]
               visited[i] = true
               stack = append(stack, i)
               for len(stack) > 0 {
                   u := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   for _, v := range adj[u] {
                       if !visited[v] && colors[v] == target {
                           visited[v] = true
                           stack = append(stack, v)
                       }
                   }
               }
           }
       }
   }
   ans := comp[0]
   if comp[1] < ans {
       ans = comp[1]
   }
   fmt.Println(ans)
}
