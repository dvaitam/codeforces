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
   deg := make([]int, n+1)
   for i := 0; i < n-1; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       deg[u]++
       deg[v]++
   }
   // if any node has degree exactly 2, configuration cannot be achieved
   for i := 1; i <= n; i++ {
       if deg[i] == 2 {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
