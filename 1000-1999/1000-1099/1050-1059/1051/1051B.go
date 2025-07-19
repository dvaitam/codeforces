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

   var l, r int64
   // Print header
   fmt.Fprintln(writer, "YES")
   // Read range
   if _, err := fmt.Fscan(reader, &l, &r); err != nil {
       return
   }
   // Output pairs
   for i := l; i <= r; i += 2 {
       fmt.Fprintf(writer, "%d %d\n", i, i+1)
   }
}
