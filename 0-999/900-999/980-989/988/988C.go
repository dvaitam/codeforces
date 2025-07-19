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

   var k int
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   // Map from exc value to first occurrence (sequence index, element index)
   type occ struct{ seq, idx int }
   occMap := make(map[int64]occ)

   for t := 1; t <= k; t++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       var sum int64
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
           sum += a[i]
       }
       for i := 0; i < n; i++ {
           exc := sum - a[i]
           if prev, found := occMap[exc]; found && prev.seq != t {
               fmt.Fprintln(writer, "YES")
               fmt.Fprintf(writer, "%d %d\n", prev.seq, prev.idx)
               fmt.Fprintf(writer, "%d %d\n", t, i+1)
               return
           } else if !found {
               occMap[exc] = occ{t, i + 1}
           }
       }
   }
   fmt.Fprintln(writer, "NO")
}
