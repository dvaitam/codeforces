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
   openLeft, openRight := 0, 0
   for i := 0; i < n; i++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       openLeft += l
       openRight += r
   }
   closeLeft := n - openLeft
   closeRight := n - openRight
   // choose minimal moves for left and right doors
   movesLeft := openLeft
   if closeLeft < movesLeft {
       movesLeft = closeLeft
   }
   movesRight := openRight
   if closeRight < movesRight {
       movesRight = closeRight
   }
   fmt.Println(movesLeft + movesRight)
}
