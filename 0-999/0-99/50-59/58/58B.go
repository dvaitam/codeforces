package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   // Build the longest chain of divisors from n down to 1
   denom := make([]int, 0, 32)
   denom = append(denom, n)
   cur := n
   tmp := n
   // Factor out primes, dividing one factor per step
   for p := 2; p*p <= tmp; p++ {
       for tmp%p == 0 {
           tmp /= p
           cur /= p
           denom = append(denom, cur)
       }
   }
   if tmp > 1 {
       // Remaining prime factor
       cur /= tmp
       denom = append(denom, cur)
       tmp = 1
   }
   // Output in decreasing order
   writer := bufio.NewWriter(os.Stdout)
   for i, v := range denom {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", v))
   }
   writer.WriteByte('\n')
   writer.Flush()
}
