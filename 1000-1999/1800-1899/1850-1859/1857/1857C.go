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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       m := n * (n - 1) / 2
       a := make([]int, m)
       for i := 0; i < m; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       k := 0
       for i := 1; i < n; i++ {
           fmt.Fprint(writer, a[k], " ")
           k += n - i
       }
       // last element
       fmt.Fprintln(writer, a[m-1])
   }
}
