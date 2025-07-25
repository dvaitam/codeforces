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

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // If no cats left, one group remains; otherwise, max groups is min(m, n-m)
   if m == 0 {
       fmt.Fprintln(out, 1)
   } else {
       rem := n - m
       if rem < m {
           fmt.Fprintln(out, rem)
       } else {
           fmt.Fprintln(out, m)
       }
   }
}
