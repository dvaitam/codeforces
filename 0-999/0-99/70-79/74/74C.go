package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b int64) int64 {
   for b != 0 {
      a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int64
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
      return
   }
   // The maximum number of non-beating balls is gcd(n-1, m-1) + 1
   result := gcd(n-1, m-1) + 1
   fmt.Println(result)
}
