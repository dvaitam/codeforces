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

   var n, k int64
   fmt.Fscan(reader, &n, &k)
   sig := k * (k + 1) / 2
   if k == 1 {
       fmt.Fprintln(writer, "YES")
       fmt.Fprintln(writer, n)
       return
   }
   if n < sig {
       fmt.Fprintln(writer, "NO")
       return
   }
   n -= sig
   q := n / k
   r := n % k
   if q > 0 || (q == 0 && r != k-1) {
       fmt.Fprintln(writer, "YES")
       for i := int64(1); i < k; i++ {
           fmt.Fprintf(writer, "%d ", i+q)
       }
       fmt.Fprintln(writer, k+q+r)
       return
   }
   if k >= 4 {
       fmt.Fprintln(writer, "YES")
       for i := int64(1); i <= k-2; i++ {
           fmt.Fprintf(writer, "%d ", i+q)
       }
       fmt.Fprintf(writer, "%d %d\n", k+q, k+q+r-1)
   } else {
       fmt.Fprintln(writer, "NO")
   }
}
