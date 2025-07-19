package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]uint64, n)
   // counts of trailing zeros
   cnt := make([]int, 64)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
       tz := bits.TrailingZeros64(s[i])
       if tz < len(cnt) {
           cnt[tz]++
       }
   }
   // find bit position with maximum count
   maxCnt := 0
   tu := -1
   for i, v := range cnt {
       if v > maxCnt {
           maxCnt = v
           tu = i
       }
   }
   // collect numbers to remove (those with tz != tu)
   var res []uint64
   for _, v := range s {
       if bits.TrailingZeros64(v) != tu {
           res = append(res, v)
       }
   }
   // output
   fmt.Fprintln(writer, len(res))
   if len(res) > 0 {
       for i, v := range res {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
