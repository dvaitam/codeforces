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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       // find last zero position (1-based)
       idx := -1
       for i, ch := range s {
           if ch == '0' {
               idx = i + 1
           }
       }
       mid := n / 2
       // output two distinct substrings of length >= floor(n/2)
       if idx > mid {
           // zero in right half
           fmt.Fprintln(writer, 1, idx, 1, idx-1)
       } else if idx < mid {
           // zero in left half or none
           fmt.Fprintln(writer, mid+1, n, mid, n-1)
       } else {
           // zero exactly at mid
           fmt.Fprintln(writer, mid+1, n, mid, n)
       }
   }
}
