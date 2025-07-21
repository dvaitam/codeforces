package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Starting pair is (1,1), if n == 1, already contains n
   if n == 1 {
       fmt.Println(0)
       return
   }
   const inf = 1000000000
   minSteps := inf
   // Try pairs (n, b) with 1 <= b < n
   for b := 1; b < n; b++ {
       x, y := n, b
       var sum int
       // Reverse Euclid: count quotients
       for y != 0 {
           sum += x / y
           x, y = y, x%y
       }
       // Only reachable if gcd(n, b) == 1
       if x != 1 {
           continue
       }
       // Steps is total additions minus one
       steps := sum - 1
       if steps < minSteps {
           minSteps = steps
       }
   }
   fmt.Println(minSteps)
}
