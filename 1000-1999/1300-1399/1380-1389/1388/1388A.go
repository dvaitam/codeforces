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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var n int
       fmt.Fscan(reader, &n)
       if n <= 30 {
           fmt.Fprintln(writer, "NO")
       } else if n == 36 {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, "5 6 10 15")
       } else if n == 40 {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, "5 6 14 15")
       } else if n == 44 {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintln(writer, "6 7 10 21")
       } else {
           fmt.Fprintln(writer, "YES")
           fmt.Fprintf(writer, "6 10 14 %d\n", n-30)
       }
   }
}
