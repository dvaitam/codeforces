package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   sBytes := make([]byte, n)
   // read string without spaces
   var s string
   fmt.Fscan(reader, &s)
   for i := 0; i < n; i++ {
       sBytes[i] = s[i]
   }
   i := 0
   for i < n-1 && k > 0 {
       // find next "47"
       if !(sBytes[i] == '4' && sBytes[i+1] == '7') {
           i++
           continue
       }
       // at position i we have "47"
       if i%2 == 0 {
           // 1-based odd -> replace with "44"
           sBytes[i+1] = '4'
       } else {
           // 1-based even -> replace with "77"
           sBytes[i] = '7'
       }
       k--
       // check for oscillation
       if i > 0 && sBytes[i-1] == '4' && sBytes[i] == '7' {
           // oscillation between i-1 and i
           if k%2 == 1 {
               // apply one more on position i-1
               if (i-1)%2 == 0 {
                   sBytes[i] = '4'
               } else {
                   sBytes[i-1] = '7'
               }
           }
           break
       }
       // move back to check overlapping
       if i > 0 {
           i--
       }
   }
   // output result
   writer := bufio.NewWriter(os.Stdout)
   writer.Write(sBytes)
   writer.Flush()
}
