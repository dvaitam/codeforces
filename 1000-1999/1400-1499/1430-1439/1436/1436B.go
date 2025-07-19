package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

// isPrime checks if x is a prime number.
func isPrime(x int) bool {
   if x < 2 {
       return false
   }
   r := int(math.Sqrt(float64(x)))
   for i := 2; i <= r; i++ {
       if x%i == 0 {
           return false
       }
   }
   return true
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       // find a: smallest i such that i+n-1 is prime and i is composite
       a := 0
       for i := 1; ; i++ {
           if isPrime(i+n-1) && !isPrime(i) {
               a = i
               break
           }
       }
       // find b: smallest i such that (n-1)*a+i is prime and i is composite
       b := 0
       for i := 1; ; i++ {
           if isPrime((n-1)*a+i) && !isPrime(i) {
               b = i
               break
           }
       }
       // output matrix
       // first n-1 rows
       for i := 1; i < n; i++ {
           for j := 1; j < n; j++ {
               fmt.Fprint(writer, 1, " ")
           }
           fmt.Fprint(writer, a, "\n")
       }
       // last row
       for i := 1; i < n; i++ {
           fmt.Fprint(writer, a, " ")
       }
       fmt.Fprint(writer, b, "\n")
   }
}
