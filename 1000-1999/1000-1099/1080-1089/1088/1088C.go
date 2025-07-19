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
   start := n * 10
   // print number of edges
   fmt.Fprintln(writer, n+1)
   // first two edges
   fmt.Fprintf(writer, "2 %d 1\n", n)
   fmt.Fprintf(writer, "1 %d %d\n", n, start)
   start--
   // remaining edges
   for i := 0; i < n-1; i++ {
       fmt.Fprintf(writer, "2 %d %d\n", i+1, start)
       start--
   }
}
