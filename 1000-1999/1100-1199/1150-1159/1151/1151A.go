package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n int
   var s string
   fmt.Fscan(rdr, &n)
   fmt.Fscan(rdr, &s)
   target := "ACTG"
   best := int(1e9)
   for i := 0; i+len(target) <= n; i++ {
       cost := 0
       for j := 0; j < len(target); j++ {
           a := int(s[i+j])
           b := int(target[j])
           d := a - b
           if d < 0 {
               d = -d
           }
           if d > 26-d {
               cost += 26 - d
           } else {
               cost += d
           }
       }
       if cost < best {
           best = cost
       }
   }
   fmt.Fprintln(w, best)
}
