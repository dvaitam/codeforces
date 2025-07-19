package main

import (
   "fmt"
)

func main() {
   var b int64
   if _, err := fmt.Scan(&b); err != nil {
       return
   }
   var cnt int64
   for i := int64(1); i*i <= b; i++ {
       if b%i == 0 {
           if b/i == i {
               cnt++
           } else {
               cnt += 2
           }
       }
   }
   fmt.Println(cnt)
}
