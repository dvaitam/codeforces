package main

import (
   "bufio"
   "fmt"
   "os"
)

func grundy(x int64) int {
   switch {
   case x <= 3:
       return 0
   case x <= 15:
       return 1
   case x <= 81:
       return 2
   case x <= 6723:
       return 0
   case x <= 50625:
       return 3
   default:
       return 1
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var x int64
   xor := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       xor ^= grundy(x)
   }
   if xor != 0 {
       fmt.Println("Furlo")
   } else {
       fmt.Println("Rublo")
   }
}
