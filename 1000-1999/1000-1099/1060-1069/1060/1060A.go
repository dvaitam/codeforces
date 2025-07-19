package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   count8 := 0
   for _, ch := range s {
       if ch == '8' {
           count8++
       }
   }
   if n >= 11 {
       maxPhones := n / 11
       if count8 < maxPhones {
           maxPhones = count8
       }
       fmt.Println(maxPhones)
   } else {
       fmt.Println(0)
   }
}
