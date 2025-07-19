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
   fmt.Println("YES")
   for i := 0; i < n; i++ {
       var a, b, c, d int
       fmt.Fscan(reader, &a, &b, &c, &d)
       if a%2 != 0 {
           if b%2 != 0 {
               fmt.Println(1)
           } else {
               fmt.Println(2)
           }
       } else {
           if b%2 != 0 {
               fmt.Println(3)
           } else {
               fmt.Println(4)
           }
       }
   }
}
