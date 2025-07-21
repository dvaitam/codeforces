package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
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

   var x, y, a, b int64
   if _, err := fmt.Fscan(reader, &x, &y, &a, &b); err != nil {
       return
   }
   g := gcd(x, y)
   lcm := x / g * y
   // count multiples of lcm in [a, b]
   if lcm > b {
       fmt.Fprintln(writer, 0)
       return
   }
   high := b / lcm
   low := (a - 1) / lcm
   count := high - low
   fmt.Fprintln(writer, count)
}
