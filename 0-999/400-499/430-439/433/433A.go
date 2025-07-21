package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   count100, count200 := 0, 0
   for i := 0; i < n; i++ {
       var w int
       if _, err := fmt.Fscan(in, &w); err != nil {
           return
       }
       if w == 100 {
           count100++
       } else if w == 200 {
           count200++
       }
   }
   total := count100*100 + count200*200
   // Total must allow half to be integer multiple of 100
   if total%200 != 0 {
       fmt.Println("NO")
    } else if count100 == 0 && count200%2 == 1 {
       fmt.Println("NO")
   } else {
       fmt.Println("YES")
   }
}
