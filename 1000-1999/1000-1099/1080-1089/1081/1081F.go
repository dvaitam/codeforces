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

   var n, t int
   // Read n and initial count t
   fmt.Fscan(reader, &n, &t)

   // TODO: implement interactive solution to reconstruct the binary sequence of length n
   // using at most 10000 queries. Each query prints "l r" and flushes, then reads
   // the updated count of 1s.

   // Placeholder: output a sequence of zeros
   result := make([]int, n)
   for i := 0; i < n; i++ {
       result[i] = 0
   }
   // Print result
   fmt.Fprint(writer, "!")
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, " %d", result[i])
   }
   fmt.Fprintln(writer)
}
