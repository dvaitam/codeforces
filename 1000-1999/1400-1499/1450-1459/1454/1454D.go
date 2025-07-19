package main

import (
   "bufio"
   "fmt"
   "os"
)

func isPrime(n uint64) bool {
   if n <= 1 {
       return false
   }
   if n <= 3 {
       return true
   }
   if n%2 == 0 || n%3 == 0 {
       return false
   }
   // check divisibility by numbers of form 6k ± 1
   for i := uint64(5); i*i <= n; i += 6 {
       if n%i == 0 || n%(i+2) == 0 {
           return false
       }
   }
   return true
}

// primeFac returns the maximum exponent and corresponding prime in n's factorization
func primeFac(n uint64) (count uint64, prime uint64) {
   m := n
   var c, p uint64
   // factor 2
   var cnt uint64
   for m%2 == 0 {
       cnt++
       m /= 2
   }
   if cnt > c {
       c = cnt
       p = 2
   }
   // factor odd primes
   // note: m is decreasing, but i*i <= m ensures correctness
   for i := uint64(3); i*i <= m; i += 2 {
       cnt = 0
       for m%i == 0 {
           cnt++
           m /= i
       }
       if cnt > c {
           c = cnt
           p = i
       }
   }
   return c, p
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n uint64
       fmt.Fscan(reader, &n)
       if isPrime(n) {
           // single prime
           fmt.Fprintln(writer, 1)
           fmt.Fprintln(writer, n)
           continue
       }
       c, p := primeFac(n)
       // output number of parts
       fmt.Fprintln(writer, c)
       rem := n
       // print c-1 copies of p
       for i := uint64(1); i < c; i++ {
           fmt.Fprint(writer, p, " ")
           rem /= p
       }
       // last part
       fmt.Fprintln(writer, rem)
   }
}
