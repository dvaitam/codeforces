package main

import "fmt"

// isPrime returns true if n is a prime number.
func isPrime(n int) bool {
   if n < 2 {
       return false
   }
   for i := 2; i*i <= n; i++ {
       if n%i == 0 {
           return false
       }
   }
   return true
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   if isPrime(n) {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
