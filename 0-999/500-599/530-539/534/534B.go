package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var v1, v2, t, d int
   fmt.Fscan(in, &v1, &v2)
   fmt.Fscan(in, &t, &d)
   total := 0
   // For each second i (0-based), speed limited by v1+i*d and v2+(t-1-i)*d
   for i := 0; i < t; i++ {
       maxFromStart := v1 + i*d
       maxFromEnd := v2 + (t-1-i)*d
       if maxFromStart < maxFromEnd {
           total += maxFromStart
       } else {
           total += maxFromEnd
       }
   }
   fmt.Fprintln(out, total)
}
