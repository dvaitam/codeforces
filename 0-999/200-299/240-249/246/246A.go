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
   if n < 3 {
       fmt.Fprintln(writer, -1)
       return
   }
   // Print counter-example: arbitrary values 70, 69, then 3..n
   fmt.Fprint(writer, 70, " ", 69)
   for i := 3; i <= n; i++ {
       fmt.Fprint(writer, " ", i)
   }
   fmt.Fprintln(writer)
}
