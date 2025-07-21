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
   half := n / 2
   sum1, sum2 := 0, 0
   for i := 0; i < n; i++ {
       ch := s[i]
       if ch != '4' && ch != '7' {
           fmt.Println("NO")
           return
       }
       d := int(ch - '0')
       if i < half {
           sum1 += d
       } else {
           sum2 += d
       }
   }
   if sum1 == sum2 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
