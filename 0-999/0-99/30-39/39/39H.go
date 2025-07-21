package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var k int
   if _, err := fmt.Fscan(in, &k); err != nil {
       return
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i := 1; i < k; i++ {
       for j := 1; j < k; j++ {
           if j > 1 {
               fmt.Fprint(w, " ")
           }
           fmt.Fprint(w, toBase(i*j, k))
       }
       fmt.Fprintln(w)
   }
}

// toBase converts a non-negative integer n to its string representation in the given base (2 <= base <= 10).
func toBase(n, base int) string {
   if n == 0 {
       return "0"
   }
   var digits []byte
   for n > 0 {
       d := n % base
       digits = append(digits, byte('0'+d))
       n /= base
   }
   // reverse
   for l, r := 0, len(digits)-1; l < r; l, r = l+1, r-1 {
       digits[l], digits[r] = digits[r], digits[l]
   }
   return string(digits)
}
