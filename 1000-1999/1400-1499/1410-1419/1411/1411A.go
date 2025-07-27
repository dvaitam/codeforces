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
   for tc := 0; tc < t; tc++ {
       var n int
       var s string
       fmt.Fscan(reader, &n, &s)
       // count consecutive ')' from end
       k := 0
       for i := n - 1; i >= 0 && s[i] == ')'; i-- {
           k++
       }
       if 2*k > n {
           fmt.Println("Yes")
       } else {
           fmt.Println("No")
       }
   }
}
