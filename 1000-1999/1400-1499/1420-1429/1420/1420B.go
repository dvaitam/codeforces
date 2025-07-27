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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       cnt := make([]int64, 32)
       for i := 0; i < n; i++ {
           var a uint64
           fmt.Fscan(reader, &a)
           if a > 0 {
               k := bits.Len64(a) - 1
               cnt[k]++
           }
       }
       var ans int64
       for _, c := range cnt {
           if c > 1 {
               ans += c * (c - 1) / 2
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
