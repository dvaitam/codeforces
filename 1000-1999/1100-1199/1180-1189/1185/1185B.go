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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   for i := 0; i < n; i++ {
       var s, t string
       fmt.Fscan(reader, &s, &t)
       ok := match(s, t)
       if ok {
           writer.WriteString("YES\n")
       } else {
           writer.WriteString("NO\n")
       }
   }
}

// match returns true if t can be produced by pressing the keys for s on a broken keyboard
func match(s, t string) bool {
   i, j := 0, 0
   n, m := len(s), len(t)
   for i < n && j < m {
       if s[i] != t[j] {
           return false
       }
       ch := s[i]
       cs, ct := 0, 0
       for i+cs < n && s[i+cs] == ch {
           cs++
       }
       for j+ct < m && t[j+ct] == ch {
           ct++
       }
       if ct < cs {
           return false
       }
       i += cs
       j += ct
   }
   // both must be fully consumed
   return i == n && j == m
}
