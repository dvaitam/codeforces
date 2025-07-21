package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // Compute suffix maximum characters
   suf := make([]byte, n+1)
   // suf[n] is zero value, smaller than any lowercase letter
   for i := n - 1; i >= 0; i-- {
       if s[i] > suf[i+1] {
           suf[i] = s[i]
       } else {
           suf[i] = suf[i+1]
       }
   }
   // Build result: take s[i] if it's the max in suffix
   res := make([]byte, 0, n)
   for i := 0; i < n; i++ {
       b := s[i]
       if b == suf[i] {
           res = append(res, b)
       }
   }
   writer.Write(res)
}
