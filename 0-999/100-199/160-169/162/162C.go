package main

import (
   "fmt"
)

func main() {
   var n int
   // Read input
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   var factors []int
   // Factorize n by trial division
   for i := 2; i*i <= n; i++ {
       for n%i == 0 {
           factors = append(factors, i)
           n /= i
       }
   }
   // If remaining n is a prime greater than 1, include it
   if n > 1 {
       factors = append(factors, n)
   }
   // Output factors separated by '*'
   for i, f := range factors {
       if i > 0 {
           fmt.Printf("*")
       }
       fmt.Printf("%d", f)
   }
}
