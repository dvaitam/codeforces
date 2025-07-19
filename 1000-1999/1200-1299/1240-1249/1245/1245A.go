package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   if b == 0 {
       return a
   }
   return gcd(b, a%b)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if gcd(a, b) != 1 {
           fmt.Fprintln(writer, "infinite")
       } else {
           fmt.Fprintln(writer, "finite")
       }
   }
}
