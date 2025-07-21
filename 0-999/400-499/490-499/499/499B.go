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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   dict := make(map[string]string, m)
   var a, b string
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &a, &b)
       dict[a] = b
   }
   // process lecture
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a)
       b = dict[a]
       if len(b) < len(a) {
           fmt.Fprint(writer, b)
       } else {
           fmt.Fprint(writer, a)
       }
       if i+1 < n {
           fmt.Fprint(writer, " ")
       }
   }
   fmt.Fprintln(writer)
}
