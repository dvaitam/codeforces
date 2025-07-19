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

   var n, k int
   // read n, k and string s
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }

   // count occurrences of uppercase letters A, B, ...
   counts := make([]int, 26)
   for i := 0; i < len(s); i++ {
       c := s[i]
       if c >= 'A' && c <= 'Z' {
           counts[c-'A']++
       }
   }

   // find minimum count among first k letters
   minx := n
   for i := 0; i < k; i++ {
       if counts[i] < minx {
           minx = counts[i]
       }
   }
   result := int64(minx) * int64(k)
   fmt.Fprint(writer, result)
}
