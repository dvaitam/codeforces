package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var w, h, a, b int64
   if _, err := fmt.Fscan(reader, &w, &h, &a, &b); err != nil {
       return
   }
   x := gcd(a, b)
   a /= x
   b /= x
   c := w / a
   d := h / b
   if c < d {
       fmt.Fprint(writer, c)
   } else {
       fmt.Fprint(writer, d)
   }
}
