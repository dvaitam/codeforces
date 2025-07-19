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
   fmt.Fscan(reader, &n)
   if n%2 == 0 {
       if n == 2 {
           fmt.Fprintln(writer, 1)
           fmt.Fprintln(writer, "1 1")
       } else {
           if n%4 == 0 {
               fmt.Fprintln(writer, 0)
           } else {
               fmt.Fprintln(writer, 1)
           }
           size := n / 2
           fmt.Fprintf(writer, "%d ", size)
           flag := true
           for i := 1; i <= n; i += 2 {
               if flag {
                   fmt.Fprintf(writer, "%d ", i)
               } else {
                   fmt.Fprintf(writer, "%d ", i+1)
               }
               flag = !flag
           }
           fmt.Fprintln(writer)
       }
   } else {
       if n == 3 {
           fmt.Fprintln(writer, 0)
           fmt.Fprintln(writer, "1 3")
       } else {
           if n%4 == 1 {
               fmt.Fprintln(writer, 1)
               fmt.Fprintf(writer, "%d ", n/2)
           } else {
               fmt.Fprintln(writer, 0)
               fmt.Fprintf(writer, "%d 1 ", n/2+1)
           }
           flag := true
           for i := 2; i <= n; i += 2 {
               if flag {
                   fmt.Fprintf(writer, "%d ", i)
               } else {
                   fmt.Fprintf(writer, "%d ", i+1)
               }
               flag = !flag
           }
           fmt.Fprintln(writer)
       }
   }
}
