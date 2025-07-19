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

   var q int
   fmt.Fscan(reader, &q)
   for ; q > 0; q-- {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       if n == 2 {
           if s[1] > s[0] {
               fmt.Fprintln(writer, "YES")
               fmt.Fprintln(writer, 2)
               fmt.Fprintf(writer, "%c %c\n", s[0], s[1])
           } else {
               fmt.Fprintln(writer, "NO")
           }
       } else {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, 2)
           fmt.Fprintf(writer, "%c %s\n", s[0], s[1:])
       }
   }
}
