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
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       if n > 1 {
           writer.WriteByte('2')
           for i := 1; i < n; i++ {
               writer.WriteByte('7')
           }
       } else {
           writer.WriteString("-1")
       }
       writer.WriteByte('\n')
   }
}
