package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   m := (n - 1) / 2
   fmt.Fprintln(writer, m)

   // print pairs
   for i := 0; i < m; i++ {
       fmt.Fprint(writer, a[m+i], " ", a[i], " ")
   }
   if n%2 == 1 {
       fmt.Fprintln(writer, a[n-1])
   } else {
       fmt.Fprint(writer, a[n-2], " ", a[n-1])
       fmt.Fprintln(writer)
   }
}
