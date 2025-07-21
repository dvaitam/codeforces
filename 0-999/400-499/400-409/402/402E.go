package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   graph := make([][]int, n)
   rev := make([][]int, n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           if x > 0 {
               graph[i] = append(graph[i], j)
               rev[j] = append(rev[j], i)
           }
       }
   }
   // check reachability from node 0
   visited := make([]bool, n)
   queue := make([]int, 0, n)
   visited[0] = true
   queue = append(queue, 0)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range graph[u] {
           if !visited[v] {
               visited[v] = true
               queue = append(queue, v)
           }
       }
   }
   for i := 0; i < n; i++ {
       if !visited[i] {
           fmt.Println("NO")
           return
       }
   }
   // check reachability in reverse graph
   for i := range visited {
       visited[i] = false
   }
   queue = queue[:0]
   visited[0] = true
   queue = append(queue, 0)
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       for _, v := range rev[u] {
           if !visited[v] {
               visited[v] = true
               queue = append(queue, v)
           }
       }
   }
   for i := 0; i < n; i++ {
       if !visited[i] {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
