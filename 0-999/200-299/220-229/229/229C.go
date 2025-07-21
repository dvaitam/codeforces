package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   deg := make([]int64, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       deg[a-1]++
       deg[b-1]++
   }
   N := int64(n)
   if n < 3 {
       fmt.Println(0)
       return
   }
   total := N * (N-1) * (N-2) / 6
   var S int64
   for i := 0; i < n; i++ {
       S += deg[i] * (N-1 - deg[i])
   }
   ans := total - S/2
   fmt.Println(ans)
}
