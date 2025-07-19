package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   flag := true
   var now int
   for i := 1; i <= n; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       if i == 1 {
           now = max(a, b)
       } else {
           if now >= a && now >= b {
               now = max(a, b)
           } else if now >= a {
               now = a
           } else if now >= b {
               now = b
           } else {
               flag = false
               break
           }
       }
   }
   if flag {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
