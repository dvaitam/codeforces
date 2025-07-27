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
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)

       hasMinus := false
       hasLeft := false
       hasRight := false
       for i := 0; i < n; i++ {
           switch s[i] {
           case '-':
               hasMinus = true
           case '<':
               hasLeft = true
           case '>':
               hasRight = true
           }
       }

       if !hasMinus {
           if !hasLeft || !hasRight {
               // all belts one-way cycle
               fmt.Fprintln(writer, n)
           } else {
               // mixed one-way, no return
               fmt.Fprintln(writer, 0)
           }
       } else {
           // count rooms adjacent to off belts
           ans := 0
           for i := 0; i < n; i++ {
               if s[i] == '-' || s[(i-1+n)%n] == '-' {
                   ans++
               }
           }
           fmt.Fprintln(writer, ans)
       }
   }
}
