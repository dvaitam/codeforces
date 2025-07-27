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
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   var score int64 = 0
   // Alice and Bob optimal play yields pairing smallest with largest
   for i := 0; i < k; i++ {
       score += int64(a[n-1-i] - a[i])
   }
   fmt.Fprintln(writer, score)
}
