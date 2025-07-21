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

   var a1, a2 int64
   if _, err := fmt.Fscan(reader, &a1, &a2); err != nil {
       return
   }
   // Output the sum of the two integers
   fmt.Fprintln(writer, a1+a2)
}
