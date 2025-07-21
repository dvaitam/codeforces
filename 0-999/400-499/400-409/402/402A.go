package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k, a, b, v int
   if _, err := fmt.Fscan(reader, &k, &a, &b, &v); err != nil {
       return
   }
   // We have unlimited boxes, each box i gets xi divisors (0 <= xi <= k-1)
   // Total divisors used <= b, sections per box = xi+1, capacity per section = v
   // For m boxes, max divisors usable = min(b, m*(k-1)), total sections = m + D
   // Total capacity = (m + D)*v
   // Find minimal m such that (m + min(b, m*(k-1)))*v >= a
   for m := 1; ; m++ {
       // total divisors assignable
       maxDiv := m * (k - 1)
       d := b
       if maxDiv < b {
           d = maxDiv
       }
       totalSections := m + d
       if totalSections * v >= a {
           fmt.Println(m)
           return
       }
   }
}
