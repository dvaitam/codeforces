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
   var A int64
   fmt.Fscan(reader, &n, &A)
   d := make([]int64, n)
   var sumd int64
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &d[i])
       sumd += d[i]
   }
   // The minimum sum of other dice is (n-1)*1
   sumMinExcl := int64(n - 1)

   // Compute for each dice the number of impossible face values
   for i := 0; i < n; i++ {
       // For face value r to be possible:
       // A - r must be between sumMinExcl and sumd - d[i]
       // => r >= A - (sumd - d[i]) and r <= A - sumMinExcl
       low := A - (sumd - d[i])
       if low < 1 {
           low = 1
       }
       high := A - sumMinExcl
       if high > d[i] {
           high = d[i]
       }
       var possible int64
       if low > high {
           possible = 0
       } else {
           possible = high - low + 1
       }
       bi := d[i] - possible
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, bi)
   }
   fmt.Fprintln(writer)
}
