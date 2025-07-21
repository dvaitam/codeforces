package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func lcm(a, b int64) int64 {
   return a / gcd(a, b) * b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var ans int64
   switch {
   case n <= 2:
       ans = n
   case n == 3:
       ans = 6
   default:
       if n%2 != 0 {
           ans = n * (n - 1) * (n - 2)
       } else {
           if n%3 != 0 {
               // pick n, n-1, n-3
               ans = lcm(lcm(n, n-1), n-3)
           } else {
               // pick n-1, n-2, n-3
               ans = (n-1) * (n-2) * (n-3)
           }
       }
   }
   fmt.Println(ans)
}
