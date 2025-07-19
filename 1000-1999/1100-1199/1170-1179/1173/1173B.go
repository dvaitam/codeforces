package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   var n int64
   if _, err := fmt.Fscan(os.Stdin, &n); err != nil {
       return
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   // number of points to output
   cnt := n/2 + 1
   fmt.Fprintln(w, cnt)
   r, c := int64(1), int64(1)
   for i := int64(0); i < n; i++ {
       fmt.Fprintf(w, "%d %d\n", r, c)
       if i%2 == 1 {
           r++
       } else {
           c++
       }
   }
}
