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
   maxH := 0
   for i := 0; i < n; i++ {
       var h int
       fmt.Fscan(reader, &h)
       if h > maxH {
           maxH = h
       }
   }
   fmt.Fprintln(writer, maxH)
}
