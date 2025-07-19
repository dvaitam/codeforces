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
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       s := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &s[i])
       }
       // check for any unique element
       ok := true
       for i := 0; i < n; i++ {
           if (i == 0 || s[i] != s[i-1]) && (i == n-1 || s[i] != s[i+1]) {
               ok = false
               break
           }
       }
       if !ok {
           fmt.Fprintln(writer, -1)
           continue
       }
       // build permutation with cyclic shifts in equal blocks
       p := make([]int, n)
       prev := 0
       for i := 0; i < n-1; i++ {
           if s[i] == s[i+1] {
               p[i] = i + 2 // shift to next index (1-based)
           } else {
               p[i] = prev + 1 // close the cycle
               prev = i + 1    // start of next block
           }
       }
       // complete last element of last block
       p[n-1] = prev + 1
       // output
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, p[i])
       }
       writer.WriteByte('\n')
   }
}
