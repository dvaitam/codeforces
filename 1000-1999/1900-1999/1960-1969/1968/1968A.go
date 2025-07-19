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
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var x int
       fmt.Fscan(reader, &x)
       best := 1
       bestVal := gcd(1, x) + 1
       // iterate from 2 to x-1 to find the optimal i
       for i := 2; i < x; i++ {
           v := gcd(i, x) + i
           if v > bestVal {
               bestVal = v
               best = i
           }
       }
       fmt.Fprintln(writer, best)
   }
}
