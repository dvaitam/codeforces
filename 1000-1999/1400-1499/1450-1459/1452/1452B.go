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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       var sum, mx int64
       for i := 0; i < n; i++ {
           var x int64
           fmt.Fscan(reader, &x)
           sum += x
           if x > mx {
               mx = x
           }
       }
       // Ensure total sum is enough to distribute blocks
       need1 := mx*int64(n-1) - sum
       if need1 < 0 {
           need1 = 0
       }
       sum += need1
       // Make sum divisible by (n-1)
       rem := sum % int64(n-1)
       var need2 int64
       if rem != 0 {
           need2 = int64(n-1) - rem
       }
       fmt.Fprintln(writer, need1+need2)
   }
}
