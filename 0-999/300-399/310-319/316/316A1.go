package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }

   // Count distinct letters A-J
   seen := make([]bool, 10)
   k := 0
   for _, c := range s {
       if c >= 'A' && c <= 'J' {
           idx := c - 'A'
           if !seen[idx] {
               seen[idx] = true
               k++
           }
       }
   }

   // Compute number of ways to assign letters
   var letterMapping int64 = 1
   first := s[0]
   if first >= 'A' && first <= 'J' {
       // First letter cannot be zero
       letterMapping = 9
       rem := k - 1
       for i := 0; i < rem; i++ {
           letterMapping *= int64(9 - i)
       }
   } else {
       // Digits or '?', letters can use any digits
       for i := 0; i < k; i++ {
           letterMapping *= int64(10 - i)
       }
   }

   // Compute number of ways for '?' positions
   var multiplier int64 = 1
   for i, c := range s {
       if c == '?' {
           if i == 0 {
               multiplier *= 9
           } else {
               multiplier *= 10
           }
       }
   }

   result := letterMapping * multiplier
   fmt.Println(result)
}
