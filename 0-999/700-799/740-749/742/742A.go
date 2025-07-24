package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var res int
   if n == 0 {
       res = 1
   } else {
       switch n % 4 {
       case 1:
           res = 8
       case 2:
           res = 4
       case 3:
           res = 2
       case 0:
           res = 6
       }
   }
   fmt.Println(res)
}
