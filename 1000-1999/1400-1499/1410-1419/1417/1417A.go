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

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, k int
       fmt.Fscan(reader, &n, &k)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       sort.Ints(a)
       mn := a[0]
       var ans int
       for i := 1; i < n; i++ {
           // number of times we can add mn without exceeding k
           ans += (k - a[i]) / mn
       }
       fmt.Fprintln(writer, ans)
   }
}
