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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }

   g := a[0]
   max := a[0]
   for i := 1; i < n; i++ {
       g = gcd(g, a[i])
       if a[i] > max {
           max = a[i]
       }
   }
   // Total numbers in final set = max/g
   total := max / g
   // Moves = numbers added = total - initial count
   moves := total - int64(n)
   if moves%2 == 1 {
       fmt.Fprint(writer, "Alice")
   } else {
       fmt.Fprint(writer, "Bob")
   }
}
