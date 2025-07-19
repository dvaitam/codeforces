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
   if n%2 == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   // line 1: 0 1 2 ... n-1
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, "%d ", i)
   }
   fmt.Fprintln(writer)
   // line 2: 1 2 ... n-1 0
   for i := 1; i < n; i++ {
       fmt.Fprintf(writer, "%d ", i)
   }
   fmt.Fprintf(writer, "0")
   fmt.Fprintln(writer)
   // line 3: odd indices then even indices
   for i := 1; i < n; i += 2 {
       fmt.Fprintf(writer, "%d ", i)
   }
   for i := 0; i < n; i += 2 {
       fmt.Fprintf(writer, "%d ", i)
   }
   fmt.Fprintln(writer)
}
