package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

// gcd computes the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // For n == 1, no tosses needed
   // Expected tosses = 0/1
   if n == 1 {
       fmt.Println("0/1")
       return
   }
   // k = floor(log2(n))
   k := bits.Len(uint(n)) - 1
   // m = number of codewords of length k: m = 2^(k+1) - n
   m := (1 << (k + 1)) - n
   // numerator = total expected tosses * n = n*(k+1) - m
   num := n*(k+1) - m
   den := n
   g := gcd(num, den)
   num /= g
   den /= g
   fmt.Printf("%d/%d", num, den)
}
