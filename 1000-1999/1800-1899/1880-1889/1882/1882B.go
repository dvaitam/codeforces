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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
      return
   }
   for ; t > 0; t-- {
      var n int
      fmt.Fscan(reader, &n)
      sets := make([]uint64, n)
      var all uint64
      for i := 0; i < n; i++ {
         var k int
         fmt.Fscan(reader, &k)
         var mask uint64
         for j := 0; j < k; j++ {
            var x int
            fmt.Fscan(reader, &x)
            mask |= 1 << uint(x-1)
         }
         sets[i] = mask
         all |= mask
      }
      ans := 0
      for e := 0; e < 50; e++ {
         if all&(1<<uint(e)) == 0 {
            continue
         }
         var u uint64
         for i := 0; i < n; i++ {
            if sets[i]&(1<<uint(e)) == 0 {
               u |= sets[i]
            }
         }
         cnt := bits.OnesCount64(u)
         if cnt > ans {
            ans = cnt
         }
      }
      fmt.Fprintln(writer, ans)
   }
}
