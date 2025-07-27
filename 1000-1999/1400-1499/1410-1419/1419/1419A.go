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
       var n int
       var s string
       fmt.Fscan(reader, &n)
       fmt.Fscan(reader, &s)
       // 1-indexed positions
       if n%2 == 1 {
           // Raze wants an odd digit at odd position
           win := false
           for i := 0; i < n; i += 2 {
               d := s[i] - '0'
               if d%2 == 1 {
                   win = true
                   break
               }
           }
           if win {
               fmt.Fprintln(writer, 1)
           } else {
               fmt.Fprintln(writer, 2)
           }
       } else {
           // Breach wants an even digit at even position
           win := false
           for i := 1; i < n; i += 2 {
               d := s[i] - '0'
               if d%2 == 0 {
                   win = true
                   break
               }
           }
           if win {
               fmt.Fprintln(writer, 2)
           } else {
               fmt.Fprintln(writer, 1)
           }
       }
   }
}
