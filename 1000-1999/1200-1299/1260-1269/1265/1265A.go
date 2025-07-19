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

   var t int
   fmt.Fscan(reader, &t)
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       b := []rune(s)
       invalid := false
       for i := 0; i < n; i++ {
           if b[i] != '?' {
               if (i > 0 && b[i] == b[i-1]) || (i+1 < n && b[i] == b[i+1]) {
                   invalid = true
                   break
               }
           }
       }
       if invalid {
           fmt.Fprintln(writer, -1)
           continue
       }
       for i := 0; i < n; i++ {
           if b[i] == '?' {
               for _, c := range []rune{'a', 'b', 'c'} {
                   if (i > 0 && b[i-1] == c) || (i+1 < n && b[i+1] == c) {
                       continue
                   }
                   b[i] = c
                   break
               }
           }
       }
       fmt.Fprintln(writer, string(b))
   }
}
