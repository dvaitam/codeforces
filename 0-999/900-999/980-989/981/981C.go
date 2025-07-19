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
   deg := make([]int, n+1)
   for i := 1; i < n; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       deg[u]++
       deg[v]++
   }
   cnt := 0
   for i := 1; i <= n; i++ {
       if deg[i] > 2 {
           cnt++
           if cnt > 1 {
               fmt.Println("No")
               return
           }
       }
   }
   // find node with maximum degree
   pre := 1
   for i := 1; i <= n; i++ {
       if deg[i] > deg[pre] {
           pre = i
       }
   }
   fmt.Println("Yes")
   fmt.Println(deg[pre])
   // connect center to all leaves
   for i := 1; i <= n; i++ {
       if i != pre && deg[i] == 1 {
           fmt.Printf("%d %d\n", pre, i)
       }
   }
}
