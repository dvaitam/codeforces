package main

import "fmt"

func main() {
   var k int
   if _, err := fmt.Scan(&k); err != nil {
       return
   }
   primes := make([]string, k)
   for i := 0; i < k; i++ {
       fmt.Scan(&primes[i])
   }
   // Output number of primes and the primes
   fmt.Print(k)
   for _, p := range primes {
       fmt.Print(" ", p)
   }
   fmt.Println()
}
