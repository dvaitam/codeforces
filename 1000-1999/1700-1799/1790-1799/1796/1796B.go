package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var a, b string
       fmt.Fscan(reader, &a, &b)
       if a == b {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, a)
           continue
       }
       if a[0] == b[0] {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintf(writer, "%c*\n", a[0])
           continue
       }
       if a[len(a)-1] == b[len(b)-1] {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintf(writer, "*%c\n", a[len(a)-1])
           continue
       }
       found := ""
       for j := 0; j+1 < len(a); j++ {
           substr := a[j : j+2]
           if strings.Contains(b, substr) {
               found = substr
               break
           }
       }
       if found == "" {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintf(writer, "*%s*\n", found)
       }
   }
}
