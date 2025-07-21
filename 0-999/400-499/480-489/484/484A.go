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

   var n int
   _, err := fmt.Fscan(reader, &n)
   if err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var l, r uint64
       fmt.Fscan(reader, &l, &r)
       ans := r
       // Try reducing one bit at a time from LSB to higher
       for k := 0; k < 63; k++ {
           if (ans>>k)&1 == 1 {
               // Clear bit k and lower bits, then set lower bits to 1
               mask := (uint64(1) << (k + 1)) - 1
               tmp := (ans & ^mask) | (mask >> 1)
               if tmp >= l {
                   ans = tmp
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
