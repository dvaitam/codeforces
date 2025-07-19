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
   for ; t > 0; t-- {
       var s string
       fmt.Fscan(reader, &s)
       n := len(s)
       // collect unique characters
       set := make(map[byte]struct{})
       for i := 0; i < n; i++ {
           set[s[i]] = struct{}{}
       }
       if len(set) == 1 {
           // all characters are same
           fmt.Fprintln(writer, string(s[0]), s[1:n-1], string(s[n-1]))
       } else {
           ok := false
           // try to pick middle part as single 'a'
           for i := 1; i < n-1; i++ {
               if s[i] == 'a' {
                   // split: [0:i), [i], (i+1:)
                   fmt.Fprintln(writer, s[:i], "a", s[i+1:])
                   ok = true
                   break
               }
           }
           if !ok {
               // default split
               fmt.Fprintln(writer, string(s[0]), s[1:n-1], string(s[n-1]))
           }
       }
   }
}
