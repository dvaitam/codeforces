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
   var s string
   fmt.Fscan(reader, &k)
   fmt.Fscan(reader, &s)

   n := len(s)
   // Special case: count substrings with zero '1's
   if k == 0 {
       var res int64
       var zeroLen int64
       for i := 0; i < n; i++ {
           if s[i] == '0' {
               zeroLen++
           } else {
               res += zeroLen * (zeroLen + 1) / 2
               zeroLen = 0
           }
       }
       res += zeroLen * (zeroLen + 1) / 2
       fmt.Fprintln(writer, res)
       return
   }

   // Record positions of '1's with sentinels
   positions := make([]int, 0, n+2)
   positions = append(positions, -1)
   for i := 0; i < n; i++ {
       if s[i] == '1' {
           positions = append(positions, i)
       }
   }
   positions = append(positions, n)
   cnt := len(positions) - 2
   if k > cnt {
       fmt.Fprintln(writer, 0)
       return
   }

   // Count substrings containing exactly k '1's
   var res int64
   for i := 1; i <= cnt-k+1; i++ {
       leftZeros := int64(positions[i] - positions[i-1])
       rightZeros := int64(positions[i+k] - positions[i+k-1])
       res += leftZeros * rightZeros
   }
   fmt.Fprintln(writer, res)
}
