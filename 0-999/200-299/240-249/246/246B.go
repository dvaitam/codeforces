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
   sum := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       sum += a
   }
   if sum%n == 0 {
       fmt.Fprintln(writer, n)
   } else {
       fmt.Fprintln(writer, n-1)
   }
}
