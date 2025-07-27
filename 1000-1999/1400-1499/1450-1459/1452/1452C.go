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
       openRound, openSquare, ans := 0, 0, 0
       for i := 0; i < len(s); i++ {
           switch s[i] {
           case '(':
               openRound++
           case '[':
               openSquare++
           case ')':
               if openRound > 0 {
                   ans++
                   openRound--
               }
           case ']':
               if openSquare > 0 {
                   ans++
                   openSquare--
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
