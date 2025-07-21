package main

import (
   "fmt"
   "os"
)

func main() {
   var a, b int64
   if _, err := fmt.Fscan(os.Stdin, &a, &b); err != nil {
       return
   }

   if a < b {
       fmt.Println(0)
       return
   }
   if a == b {
       fmt.Println("infinity")
       return
   }
   d := a - b
   var cnt int64
   for i := int64(1); i*i <= d; i++ {
       if d%i == 0 {
           x := i
           y := d / i
           if x > b {
               cnt++
           }
           if y != x && y > b {
               cnt++
           }
       }
   }
   fmt.Println(cnt)
}
