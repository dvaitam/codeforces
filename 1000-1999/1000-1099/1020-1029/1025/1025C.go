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
   n := len(s)
   // Duplicate to handle circularity
   s2 := s + s
   mx, cur := 0, 0
   var prev byte
   for i := 0; i < len(s2); i++ {
       c := s2[i]
       if i == 0 || c != prev {
           cur++
       } else {
           cur = 1
       }
       prev = c
       if cur > n {
           cur = n
       }
       if cur > mx {
           mx = cur
       }
   }
   fmt.Println(mx)
}
