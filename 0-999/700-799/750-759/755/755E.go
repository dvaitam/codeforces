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
   if n < 4 {
       fmt.Fprintln(writer, -1)
       return
   }
   if n == 4 && m != 3 {
       fmt.Fprintln(writer, -1)
       return
   }
   if m == 2 {
       fmt.Fprintln(writer, n-1)
       for i := 1; i <= n-1; i++ {
           fmt.Fprintf(writer, "%d %d\n", i, i+1)
       }
       return
   }
   if m == 3 {
       fmt.Fprintln(writer, 3+(n-4)*2)
       for i := 1; i <= 3; i++ {
           fmt.Fprintf(writer, "%d %d\n", i, i+1)
       }
       for i := 5; i <= n; i++ {
           fmt.Fprintf(writer, "%d %d\n", i, 1)
           fmt.Fprintf(writer, "%d %d\n", i, 2)
       }
       return
   }
   fmt.Fprintln(writer, -1)
}
