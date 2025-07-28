package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of a and b using the Euclidean algorithm.
func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   fmt.Println(gcd(a, b))
}
