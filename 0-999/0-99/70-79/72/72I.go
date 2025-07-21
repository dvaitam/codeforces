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
   s := fmt.Sprint(n)
   total := len(s)
   divCount := 0
   for i := 0; i < len(s); i++ {
       d := s[i] - '0'
       if d != 0 && n%int(d) == 0 {
           divCount++
       }
   }
   switch {
   case divCount == 0:
       fmt.Println("upset")
   case divCount == total:
       fmt.Println("happier")
   default:
       fmt.Println("happy")
   }
}
