package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Build result
   res := make([]rune, len(s))
   for i, ch := range s {
       d := ch - '0'
       inv := '0' + (9 - d)
       if i == 0 {
           // First digit: avoid leading zero
           if ch == '9' {
               res[i] = ch
           } else if d > 4 {
               // safe to invert, won't be '0'
               res[i] = inv
           } else {
               res[i] = ch
           }
       } else {
           // Other digits: invert if larger than 4
           if d > 4 {
               res[i] = inv
           } else {
               res[i] = ch
           }
       }
   }
   fmt.Println(string(res))
}
