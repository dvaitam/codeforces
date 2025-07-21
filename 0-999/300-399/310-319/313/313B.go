package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
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
   // prefix sums of equal adjacent characters
   p := make([]int, n+1)
   for i := 1; i < n; i++ {
       p[i+1] = p[i]
       if s[i] == s[i-1] {
           p[i+1]++
       }
   }

   var m int
   if _, err := fmt.Fscan(reader, &m); err != nil {
       return
   }
   for i := 0; i < m; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       ans := p[r] - p[l]
       writer.WriteString(strconv.Itoa(ans))
       writer.WriteByte('\n')
   }
}
