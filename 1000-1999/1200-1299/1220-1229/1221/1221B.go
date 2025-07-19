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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if (i+j)%2 == 0 {
               writer.WriteByte('W')
           } else {
               writer.WriteByte('B')
           }
       }
       writer.WriteByte('\n')
   }
}
