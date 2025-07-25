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

   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   diff := b - a
   if diff < 0 {
       diff = -diff
   }
   var bestK int64 = 0
   // initialize best LCM to a large value
   var bestLCM uint64 = ^uint64(0)
   // iterate over divisors of diff
   for d := int64(1); d*d <= diff; d++ {
       if diff%d != 0 {
           continue
       }
       // check divisor d
       k := (d - a%d) % d
       A := a + k
       B := b + k
       // gcd(A,B) == d
       // lcm = A/d * B
       lcm := uint64(A/d) * uint64(B)
       if lcm < bestLCM || (lcm == bestLCM && k < bestK) {
           bestLCM = lcm
           bestK = k
       }
       // check paired divisor
       d2 := diff / d
       if d2 != d {
           k2 := (d2 - a%d2) % d2
           A2 := a + k2
           B2 := b + k2
           lcm2 := uint64(A2/d2) * uint64(B2)
           if lcm2 < bestLCM || (lcm2 == bestLCM && k2 < bestK) {
               bestLCM = lcm2
               bestK = k2
           }
       }
   }
   // if a == b, diff == 0, bestK remains 0
   fmt.Fprintln(writer, bestK)
}
