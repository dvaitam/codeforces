package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n <= 0 {
       return
   }
   var sb strings.Builder
   sb.WriteByte('a')
   n--
   z := byte('a')
   // append pairs
   for i := 1; i <= n/2; i++ {
       if i%2 == 1 {
           z = 'b'
           sb.WriteString("bb")
       } else {
           z = 'a'
           sb.WriteString("aa")
       }
   }
   // handle odd remainder
   if n%2 == 1 {
       if z == 'a' {
           sb.WriteByte('b')
       } else {
           sb.WriteByte('a')
       }
   }
   // output result
   fmt.Print(sb.String())
}
