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
   // diff[i] = total give - total take for friend i
   diff := make([]int64, n+1)
   for i := 0; i < m; i++ {
       var a, b int
       var c int64
       fmt.Fscan(reader, &a, &b, &c)
       diff[a] += c
       diff[b] -= c
   }
   var ans int64
   for i := 1; i <= n; i++ {
       if diff[i] > 0 {
           ans += diff[i]
       }
   }
   fmt.Println(ans)
}
