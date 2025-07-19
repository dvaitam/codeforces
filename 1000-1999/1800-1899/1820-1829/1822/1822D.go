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
       if n == 1 {
           fmt.Fprintln(writer, 1)
           continue
       }
       if n%2 != 0 {
           fmt.Fprintln(writer, -1)
           continue
       }
       // build permutation
       // start with n
       // then pairs (n-1,2), (n-3,4), ..., then 1
       // print with spaces
       // use fmt.Fprint for speed
       fmt.Fprint(writer, n)
       fmt.Fprint(writer, " ")
       for i, j := n-1, 2; j <= n-2; i, j = i-2, j+2 {
           fmt.Fprint(writer, i)
           fmt.Fprint(writer, " ")
           fmt.Fprint(writer, j)
           fmt.Fprint(writer, " ")
       }
       fmt.Fprintln(writer, 1)
   }
}
