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
   // discard edges
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // stub: output zeros for each x
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString("0")
   }
   writer.WriteByte('\n')
}
