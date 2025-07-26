package main

import (
   "bufio"
   "fmt"
   "os"
)

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
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t, n int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       fmt.Fscan(in, &n)
       if isPrime(n) {
           fmt.Fprintln(out, "YES")
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
