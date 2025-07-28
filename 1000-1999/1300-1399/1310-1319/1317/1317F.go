package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int) int {
   if a < 0 {
       a = -a
   }
   if b < 0 {
       b = -b
   }
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var a, b int
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   fmt.Fprint(writer, gcd(a, b))
}
