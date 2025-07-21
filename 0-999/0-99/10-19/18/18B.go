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

   var n, d, m, l int64
   if _, err := fmt.Fscan(in, &n, &d, &m, &l); err != nil {
       return
   }
   limit := n * m
   x := d
   for {
       if x >= limit || x%m > l {
           fmt.Fprintln(out, x)
           return
       }
       x += d
   }
}
