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

   var n, m int
   var s string
   fmt.Fscan(reader, &n, &m, &s)
   // compute prefix-function (KMP) for s
   pi := make([]int, n)
   for i := 1; i < n; i++ {
       j := pi[i-1]
       for j > 0 && s[i] != s[j] {
           j = pi[j-1]
       }
       if s[i] == s[j] {
           j++
       }
       pi[i] = j
   }
   // length of minimal suffix to append
   l := pi[n-1]
   p := s[l:]
   // output
   writer.WriteString(s)
   for i := 1; i < m; i++ {
       writer.WriteString(p)
   }
}
