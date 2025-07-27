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
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tt := 0; tt < t; tt++ {
       var n int
       var s string
       fmt.Fscan(reader, &n)
       fmt.Fscan(reader, &s)
       ok := false
       // try removing one contiguous substring to get "2020"
       // check prefix of length k and suffix of length 4-k
       for k := 0; k <= 4; k++ {
           // prefix [0:k)
           // suffix [n-(4-k):n)
           if k <= n && 4-k <= n {
               if s[:k] + s[n-(4-k):] == "2020" {
                   ok = true
                   break
               }
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
