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
   ar := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &ar[i])
   }

   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       fmt.Fprintf(writer, "%d %d\n", i+1, i+1)
   }
   for i := 0; i < n; i++ {
       sel := ar[i] / 2
       if i >= sel {
           fmt.Fprintf(writer, "%d ", i-sel+1)
       } else if i+sel < n {
           fmt.Fprintf(writer, "%d ", i+sel+1)
       } else {
           // invalid case
           panic("invalid selection")
       }
   }
   fmt.Fprintln(writer)
}
