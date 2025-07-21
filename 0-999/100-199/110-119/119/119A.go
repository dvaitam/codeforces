package main

import (
   "bufio"
   "fmt"
   "os"
)

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var a, b, n int
   if _, err := fmt.Fscan(in, &a, &b, &n); err != nil {
       return
   }
   turn := 0 // 0 for Simon, 1 for Antisimon
   for {
       // determine stones to take
       var take int
       if turn == 0 {
           take = gcd(a, n)
       } else {
           take = gcd(b, n)
       }
       // check if current player can take
       if n < take {
           // current player loses
           if turn == 0 {
               fmt.Println(1)
           } else {
               fmt.Println(0)
           }
           return
       }
       n -= take
       turn ^= 1
   }
}
