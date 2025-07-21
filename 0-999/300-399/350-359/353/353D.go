package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   // Read input string of M and F without spaces
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   var m, t int64
   // m: count of boys seen so far (M)
   // t: current time needed
   for i := 0; i < len(s); i++ {
       if s[i] == 'M' {
           m++
       } else {
           // s[i] == 'F'
           if m > 0 {
               // This girl needs to swap past all previous boys
               // But cannot overtake previous girls: ensure at least t+1
               if t+1 > m {
                   t++
               } else {
                   t = m
               }
           }
       }
   }
   fmt.Println(t)
}
