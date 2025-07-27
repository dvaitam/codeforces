package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   // Precompute primes up to maxN
   const maxN = 1000000
   isPrime := make([]bool, maxN+1)
   for i := 2; i <= maxN; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= maxN; i++ {
       if isPrime[i] {
           for j := i * i; j <= maxN; j += i {
               isPrime[j] = false
           }
       }
   }
   pi := make([]int, maxN+1)
   for i := 1; i <= maxN; i++ {
       pi[i] = pi[i-1]
       if isPrime[i] {
           pi[i]++
       }
   }

   for i := 0; i < t; i++ {
       var n int
       fmt.Fscan(reader, &n)
       // lonely count = 1 (for 1) + primes p such that p*p > n
       // = 1 + (pi[n] - pi[floor(sqrt(n))])
       s := int(math.Sqrt(float64(n)))
       if (s+1)*(s+1) <= n {
           s++
       }
       for s*s > n {
           s--
       }
       lonely := 1 + (pi[n] - pi[s])
       // for n = 0, though not possible, correct: lonely = 0
       if n == 0 {
           lonely = 0
       }
       fmt.Fprintln(writer, lonely)
   }
}
