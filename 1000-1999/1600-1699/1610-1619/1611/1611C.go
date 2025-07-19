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

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       if a[0] != n && a[n-1] != n {
           fmt.Fprintln(writer, -1)
       } else if a[0] == n {
           for i := n - 1; i >= 0; i-- {
               fmt.Fprint(writer, a[i], " ")
           }
           fmt.Fprintln(writer)
       } else {
           for i := n - 2; i >= 0; i-- {
               fmt.Fprint(writer, a[i], " ")
           }
           fmt.Fprintln(writer, n)
       }
   }
}
