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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   cnt := 0
   for i := 0; i < m; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if a[x] == 0 {
           cnt++
           a[x]++
           if cnt >= n {
               writer.WriteByte('1')
               for j := 1; j <= n; j++ {
                   a[j]--
                   if a[j] == 0 {
                       cnt--
                   }
               }
           } else {
               writer.WriteByte('0')
           }
       } else {
           writer.WriteByte('0')
           a[x]++
       }
   }
}
