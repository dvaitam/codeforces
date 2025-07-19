package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       var s string
       fmt.Fscan(reader, &s)
       if n < 2 {
           fmt.Println(-1)
           continue
       }
       var fi, se byte
       ct := 0
       for i := 0; i < n; i++ {
           c := s[i]
           if (c-'0')%2 == 1 {
               if fi == 0 {
                   fi = c
                   ct++
               } else if se == 0 {
                   se = c
                   ct++
               }
           }
           if ct == 2 {
               break
           }
       }
       if ct < 2 {
           fmt.Println(-1)
       } else {
           fmt.Printf("%c%c\n", fi, se)
       }
   }
}
