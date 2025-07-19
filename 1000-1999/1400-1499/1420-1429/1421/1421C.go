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
   pal := true
   for i := 0; i < n/2; i++ {
       if s[i] != s[n-1-i] {
           pal = false
           break
       }
   }
   if pal {
       fmt.Println(0)
   } else {
       fmt.Println(3)
       fmt.Printf("L %d\n", n-1)
       fmt.Printf("R %d\n", n-1)
       fmt.Printf("R %d\n", 2*n-1)
   }
}
