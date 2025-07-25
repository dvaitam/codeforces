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

   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   if m >= n && k >= n {
       fmt.Fprintln(out, "Yes")
   } else {
       fmt.Fprintln(out, "No")
   }
}
