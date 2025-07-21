package main

import (
   "fmt"
)

// isPrime checks if n is a prime number
func isPrime(n int) bool {
   if n < 2 {
      return false
   }
   if n%2 == 0 {
      return n == 2
   }
   if n%3 == 0 {
      return n == 3
   }
   for i := 5; i*i <= n; i += 6 {
      if n%i == 0 || n%(i+2) == 0 {
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
   // check no zero digits
   temp := n
   for temp > 0 {
      if temp%10 == 0 {
         fmt.Println("NO")
         return
      }
      temp /= 10
   }
   // check n is prime
   if !isPrime(n) {
      fmt.Println("NO")
      return
   }
   // compute highest power of 10 less than n
   p := 1
   temp = n
   for temp > 0 {
      p *= 10
      temp /= 10
   }
   p /= 10
   // check all suffixes
   for p >= 10 {
      suffix := n % p
      if !isPrime(suffix) {
         fmt.Println("NO")
         return
      }
      p /= 10
   }
   fmt.Println("YES")
}
