package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd computes the greatest common divisor of x and y using the Euclidean algorithm.
func gcd(x, y int64) int64 {
   if x < 0 {
       x = -x
   }
   if y < 0 {
       y = -y
   }
   for y != 0 {
       x, y = y, x%y
   }
   return x
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a, b int64
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   result := gcd(a, b)
   fmt.Println(result)
}
