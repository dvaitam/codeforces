package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   // result slice
   var c []rune
   var k int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s)
       if i == 0 {
           k = len(s)
           c = make([]rune, k)
           for j := 0; j < k; j++ {
               c[j] = '?'
           }
       }
       // process string s
       for j, ch := range s {
           if j >= k {
               break
           }
           if ch == '?' {
               continue
           }
           if c[j] == '?' {
               c[j] = ch
           } else if c[j] != ch {
               c[j] = 'T'
           }
       }
   }
   // finalize output
   for j := 0; j < k; j++ {
       switch c[j] {
       case '?':
           c[j] = 'x'
       case 'T':
           c[j] = '?'
       }
   }
   // write result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   writer.WriteString(string(c))
}
