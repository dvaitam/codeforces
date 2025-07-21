package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   counts := make([]int64, 33)
   for i := 0; i < n; i++ {
       var x uint
       fmt.Fscan(reader, &x)
       pc := bits.OnesCount(x)
       counts[pc]++
   }
   var ans int64
   for _, c := range counts {
       if c > 1 {
           ans += c * (c - 1) / 2
       }
   }
   fmt.Fprintln(writer, ans)
}
