package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   var current int64 = 1
   var total int64 = 0
   for i := int64(0); i < m; i++ {
       var a int64
       fmt.Fscan(reader, &a)
       if a >= current {
           total += a - current
       } else {
           total += n - (current - a)
       }
       current = a
   }
   fmt.Println(total)
}
