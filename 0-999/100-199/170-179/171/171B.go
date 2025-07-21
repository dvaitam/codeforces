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

   var a int64
   if _, err := fmt.Fscan(reader, &a); err != nil {
       return
   }
   // Compute sum of first a natural numbers: 1 + 2 + ... + a
   result := a * (a + 1) / 2
   fmt.Fprintln(writer, result)
}
