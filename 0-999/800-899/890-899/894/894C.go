package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   if n == 0 {
       fmt.Fprintln(writer, "-1")
       return
   }
   t := a[0]
   for i := 1; i < n; i++ {
       t = gcd(a[i], t)
   }
   if t != a[0] {
       fmt.Fprintln(writer, "-1")
       return
   }
   total := 2*int64(n) - 1
   fmt.Fprintln(writer, total)
   // print sequence
   fmt.Fprint(writer, a[0])
   for i := 1; i < n; i++ {
       fmt.Fprintf(writer, " %d %d", a[i], a[0])
   }
}
