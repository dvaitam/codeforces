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

   var n, k, d int
   if _, err := fmt.Fscan(reader, &n, &k, &d); err != nil {
       return
   }
   // Compute smallest power b = k^i such that b > n or i reaches d
   b := int64(1)
   var i int
   for i = 0; i < d; i++ {
       b *= int64(k)
       if b > int64(n) {
           break
       }
   }
   // If k^d < n, impossible
   if i == d && b < int64(n) {
       fmt.Fprintln(writer, -1)
       return
   }
   // Prepare initial block size
   b /= int64(k)
   // Generate d sequences
   for dd := 0; dd < d; dd++ {
       if b > 0 {
           for j := 0; j < n; j++ {
               val := int((int64(j)/b)%int64(k)) + 1
               fmt.Fprint(writer, val, " ")
           }
           fmt.Fprintln(writer)
           b /= int64(k)
       } else {
           // Remaining rows filled with 1
           for j := 0; j < n; j++ {
               fmt.Fprint(writer, "1 ")
           }
           fmt.Fprintln(writer)
       }
   }
}
