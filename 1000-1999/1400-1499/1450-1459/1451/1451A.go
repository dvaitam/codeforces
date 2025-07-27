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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n int64
       fmt.Fscan(reader, &n)
       var ans int
       switch {
       case n == 1:
           ans = 0
       case n == 2:
           ans = 1
       case n == 3:
           ans = 2
       case n%2 == 0:
           ans = 2
       default:
           ans = 3
       }
       fmt.Fprintln(writer, ans)
   }
}
