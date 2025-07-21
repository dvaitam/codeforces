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
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   g := 0
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if i == 0 {
           g = x
       } else {
           g = gcd(g, x)
       }
   }
   fmt.Println(g * n)
}
