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
       var x, k int
       fmt.Fscan(reader, &x, &k)
       if x%k != 0 {
           fmt.Fprintln(writer, 1)
           fmt.Fprintln(writer, x)
       } else {
           fmt.Fprintln(writer, 2)
           fmt.Fprintln(writer, x+1, -1)
       }
   }
}
