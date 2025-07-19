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
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       // compute maximum time
       var ans int64
       for x := 3; x <= n+1; x++ {
           d := int64(n - x/2)
           ans += d * d
       }
       // output answer
       fmt.Fprintln(writer, ans)
       // construct initial permutation p
       half := n / 2
       // print p
       // first element
       fmt.Fprintf(writer, "%d", half+1)
       // next elements from 1 to half-1
       for x := 1; x < half; x++ {
           fmt.Fprintf(writer, " %d", x)
       }
       // elements from half+2 to n
       for x := half + 2; x <= n; x++ {
           fmt.Fprintf(writer, " %d", x)
       }
       // last element
       if n >= 2 {
           fmt.Fprintf(writer, " %d", half)
       }
       fmt.Fprintln(writer)
       // number of operations
       m := n - 1
       fmt.Fprintln(writer, m)
       // operations: for x from half down to 2, swap x and n
       for x := half; x >= 2; x-- {
           fmt.Fprintf(writer, "%d %d\n", x, n)
       }
       // operations: for x from half+1 to n, swap x and 1
       for x := half + 1; x <= n; x++ {
           fmt.Fprintf(writer, "%d %d\n", x, 1)
       }
   }
}
