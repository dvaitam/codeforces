package main

import (
   "bufio"
   "fmt"
   "os"
)

// compute totient of n using its prime factorization
func totient(n int) int {
   result := n
   m := n
   for i := 2; i*i <= m; i++ {
       if m%i == 0 {
           // i is a prime factor
           for m%i == 0 {
               m /= i
           }
           result = result / i * (i - 1)
       }
   }
   if m > 1 {
       // m is a prime
       result = result / m * (m - 1)
   }
   return result
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var p int
   if _, err := fmt.Fscan(reader, &p); err != nil {
       return
   }
   // number of primitive roots modulo p is phi(phi(p)) = phi(p-1)
   n := p - 1
   fmt.Println(totient(n))
}
