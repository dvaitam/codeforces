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
   f := make([]string, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       f[i] = s
   }
   visited := make([]bool, n)
   queue := make([]int, 0, n)
   // person 1 (index 0) sends initially
   visited[0] = true
   queue = append(queue, 0)
   count := make([]int, n)
   for head := 0; head < len(queue); head++ {
       u := queue[head]
       for v := 0; v < n; v++ {
           if f[u][v] == '1' {
               count[v]++
               // enqueue for forwarding if first receipt and not person n
               if v != n-1 && !visited[v] {
                   visited[v] = true
                   queue = append(queue, v)
               }
           }
       }
   }
   // output number of copies received by person n (index n-1)
   fmt.Println(count[n-1])
}
