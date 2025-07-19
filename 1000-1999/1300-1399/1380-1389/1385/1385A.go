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
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       switch {
       case a == b && b == c:
           fmt.Fprintln(writer, "YES")
           fmt.Fprintf(writer, "%d %d %d\n", a, b, c)
       case a == b || b == c || c == a:
           maxi := a
           if b > maxi {
               maxi = b
           }
           if c > maxi {
               maxi = c
           }
           if a == b {
               if a == maxi {
                   fmt.Fprintln(writer, "YES")
                   fmt.Fprintf(writer, "%d %d %d\n", maxi, 1, c)
               } else {
                   fmt.Fprintln(writer, "NO")
               }
           } else if b == c {
               if b == maxi {
                   fmt.Fprintln(writer, "YES")
                   fmt.Fprintf(writer, "%d %d %d\n", a, maxi, 1)
               } else {
                   fmt.Fprintln(writer, "NO")
               }
           } else {
               // a == c
               if a == maxi {
                   fmt.Fprintln(writer, "YES")
                   fmt.Fprintf(writer, "%d %d %d\n", 1, b, maxi)
               } else {
                   fmt.Fprintln(writer, "NO")
               }
           }
       default:
           fmt.Fprintln(writer, "NO")
       }
   }
}
