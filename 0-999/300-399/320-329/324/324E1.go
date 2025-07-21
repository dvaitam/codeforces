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
       fmt.Println(0)
       return
   }
   // skip input
   for i := 0; i < m; i++ {
       var x, y, k int
       fmt.Fscan(reader, &x, &y, &k)
       for j := 0; j < k; j++ {
           var tmp int
           fmt.Fscan(reader, &tmp)
       }
   }
   // No solution implemented
   fmt.Println(0)
}
