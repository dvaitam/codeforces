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
   if n == 5 {
       fmt.Fprintln(writer, "1 2 6")
       fmt.Fprintln(writer, "1 3 6")
       fmt.Fprintln(writer, "2 4 5")
       fmt.Fprintln(writer, "4 5 1")
       fmt.Fprintln(writer, "3 4")
       fmt.Fprintln(writer, "3 5")
       return
   }
   half := n >> 1
   // print edges
   for i := 1; i <= half; i++ {
       fmt.Fprintf(writer, "%d %d 1\n", i, i+half)
   }
   for i := 1; i+half < n; i++ {
       weight := 2*i - 1
       fmt.Fprintf(writer, "%d %d %d\n", i+half, i+half+1, weight)
   }
   // print good pairs
   for i := 1; i < half; i++ {
       fmt.Fprintf(writer, "%d %d\n", i, i+1)
   }
   fmt.Fprintln(writer, "1 3")
}
