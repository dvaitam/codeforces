package main

import (
   "fmt"
)

func main() {
   var n, k int
   if _, err := fmt.Scan(&n, &k); err != nil {
       return
   }
   // Sieve of Eratosthenes up to n
   isPrime := make([]bool, n+1)
   for i := 2; i <= n; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= n; i++ {
       if isPrime[i] {
           for j := i * i; j <= n; j += i {
               isPrime[j] = false
           }
       }
   }
   // Collect primes up to n
   primes := make([]int, 0)
   for i := 2; i <= n; i++ {
       if isPrime[i] {
           primes = append(primes, i)
       }
   }
   // Count primes of the form p_i + p_{i+1} + 1
   count := 0
   for i := 0; i+1 < len(primes); i++ {
       s := primes[i] + primes[i+1] + 1
       if s <= n && isPrime[s] {
           count++
       }
   }
   // Check against k
   if count >= k {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
