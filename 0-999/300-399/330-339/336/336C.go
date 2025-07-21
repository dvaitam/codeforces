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
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // Try bits from high to low to maximize beauty (trailing zeros)
   for j := 31; j >= 0; j-- {
       var subset []int64
       for _, v := range a {
           if (v>>j)&1 == 1 {
               subset = append(subset, v)
           }
       }
       if len(subset) == 0 {
           continue
       }
       // Compute bitwise AND of subset
       common := subset[0]
       for k := 1; k < len(subset); k++ {
           common &= subset[k]
       }
       // Count trailing zeros of common
       tz := 0
       for common&1 == 0 {
           tz++
           common >>= 1
       }
       if tz == j {
           fmt.Fprintln(writer, len(subset))
           for i, v := range subset {
               if i > 0 {
                   fmt.Fprint(writer, " ")
               }
               fmt.Fprint(writer, v)
           }
           fmt.Fprintln(writer)
           return
       }
   }
   // Fallback: print first element
   fmt.Fprintln(writer, 1)
   fmt.Fprintln(writer, a[0])
}
