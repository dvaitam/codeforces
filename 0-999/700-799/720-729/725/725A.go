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
   s := make([]byte, n)
   // read string (skip whitespace handling)
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // count leading '<'
   left := 0
   for i := 0; i < n; i++ {
       if s[i] == '<' {
           left++
       } else {
           break
       }
   }
   // count trailing '>'
   right := 0
   for i := n - 1; i >= 0; i-- {
       if s[i] == '>' {
           right++
       } else {
           break
       }
   }
   fmt.Println(left + right)
}
