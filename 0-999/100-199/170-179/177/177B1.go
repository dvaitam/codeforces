package main

import (
   "fmt"
)

// smallestPrime returns the smallest prime factor of n (n>=2).
func smallestPrime(n int64) int64 {
   if n%2 == 0 {
       return 2
   }
   for i := int64(3); i*i <= n; i += 2 {
       if n%i == 0 {
           return i
       }
   }
   return n // n is prime
}

func main() {
   var n int64
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var sum int64
   for n > 1 {
       sum += n
       p := smallestPrime(n)
       n = n / p
   }
   sum += 1 // add final state when n == 1
   fmt.Println(sum)
}
