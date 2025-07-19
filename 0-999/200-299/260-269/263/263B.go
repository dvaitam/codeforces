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

   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Fscan(reader, &a[i]); err != nil {
           return
       }
   }
   sort.Ints(a)
   if k > n {
       fmt.Fprint(writer, -1)
   } else {
       fmt.Fprintf(writer, "%d 0", a[n-k])
   }
}
