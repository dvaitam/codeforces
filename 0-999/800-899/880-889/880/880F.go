package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if isPrime(n) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}

// isPrime checks whether n is a prime number using trial division.
func isPrime(n int64) bool {
   if n < 2 {
       return false
   }
   if n%2 == 0 {
       return n == 2
   }
   for i := int64(3); i*i <= n; i += 2 {
       if n%i == 0 {
           return false
       }
   }
   return true
}
