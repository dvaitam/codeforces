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

   var N, S int
   fmt.Fscan(reader, &N, &S)

   biggest := S - N + 1
   if S - biggest + 1 >= biggest {
       fmt.Fprintln(writer, "NO")
   } else {
       fmt.Fprintln(writer, "YES")
       fmt.Fprint(writer, biggest)
       for i := 0; i < N-1; i++ {
           fmt.Fprint(writer, " ", 1)
       }
       fmt.Fprintln(writer)
       fmt.Fprintln(writer, S - biggest + 1)
   }
}
